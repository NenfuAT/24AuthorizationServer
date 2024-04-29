package controller

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/NenfuAT/24AuthorizationServer/model"
	"github.com/NenfuAT/24AuthorizationServer/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Ginによる認可エンドポイント
func Auth(c *gin.Context) {
	query := c.Request.URL.Query()
	requiredParameter := []string{"response_type", "client_id", "redirect_uri"}

	// 必須パラメータのチェック
	for _, v := range requiredParameter {
		if !query.Has(v) {
			log.Printf("%s is missing", v)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid_request. %s is missing", v),
			})
			return
		}
	}

	// client_idの確認
	if model.ClientInfo.Id != query.Get("client_id") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "client_id is not match",
		})
		return
	}

	// レスポンスタイプはコードのみ
	if query.Get("response_type") != "code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "only support code",
		})
		return
	}

	// セッションの保存
	sessionId := uuid.New().String()
	session := model.Session{
		Client:              query.Get("client_id"),
		State:               query.Get("state"),
		Scopes:              query.Get("scope"),
		RedirectUri:         query.Get("redirect_uri"),
		CodeChallenge:       query.Get("code_challenge"),
		CodeChallengeMethod: query.Get("code_challenge_method"),
	}
	// scopeの確認、OAuthかOIDCか
	// 組み合わせへの対応は面倒なので "openid profile" で固定
	if query.Get("scope") == "openid profile" {
		session.Oidc = true
	} else {
		session.CodeChallenge = query.Get("code_challenge")
		session.CodeChallengeMethod = query.Get("code_challenge_method")
	}
	model.SessionList[sessionId] = session
	// セッションIDをCookieにセット
	cookie := &http.Cookie{
		Name:  "session",
		Value: sessionId,
		Path:  "/",
	}
	http.SetCookie(c.Writer, cookie)

	// ログイン&権限認可の画面を表示
	err := model.Templates["login"].Execute(c.Writer, struct {
		ClientId string
		Scope    string
	}{
		ClientId: session.Client,
		Scope:    session.Scopes,
	})
	if err != nil {
		log.Println("Error rendering login template:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
		return
	}

	log.Println("Returned login page...")
}

// Ginによるログイン認証
func AuthCheck(c *gin.Context) {
	loginUser := c.PostForm("email")
	password := c.PostForm("password")

	user, err := model.GetUserByEmailAndPassword(loginUser, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "login failed",
		})
		return
	}

	cookie, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "authsession not found",
		})
		return
	}

	v, exists := model.SessionList[cookie]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid session",
		})
		return
	}

	// 認可コードの生成と保存
	authCodeString := uuid.New().String()
	authData := model.AuthCode{
		User:        user.Email,
		ClientId:    v.Client,
		Scopes:      v.Scopes,
		RedirectUri: v.RedirectUri,
		ExpiresAt:   time.Now().Unix() + 300, // 5分後に有効期限が切れる
	}

	model.AuthCodeList[authCodeString] = authData

	log.Printf("Auth code accepted: %s\n", authCodeString)

	// リダイレクトの設定
	location := fmt.Sprintf("%s?code=%s&state=%s", v.RedirectUri, authCodeString, v.State)
	c.Redirect(http.StatusFound, location) // 302リダイレクト
}

// トークンエンドポイント
func TokenHandler(c *gin.Context) {
	//Cookieからセッションを取得
	cookie, err := c.Cookie("session")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session not found",
		})
		return
	}

	// 必須パラメータの確認
	requiredParameters := []string{"grant_type", "code", "client_id", "redirect_uri"}
	for _, param := range requiredParameters {
		if c.PostForm(param) == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("invalid_request. %s is missing", param),
			})
			return
		}
	}

	// グラントタイプの確認
	if c.PostForm("grant_type") != "authorization_code" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request. not support type",
		})
		return
	}

	// 認可コードの取得と確認
	authCode, exists := model.AuthCodeList[c.PostForm("code")]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request. authorization code not found",
		})
		return
	}

	// クライアントIDとリダイレクトURIのチェック
	clientID := c.PostForm("client_id")
	if authCode.ClientId != clientID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request. client_id does not match",
		})
		return
	}

	redirectURI := c.PostForm("redirect_uri")
	if authCode.RedirectUri != redirectURI {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request. redirect_uri does not match",
		})
		return
	}

	// 認可コードの有効期限確認
	if authCode.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request. auth code has expired",
		})
		return
	}

	// PKCEの確認
	session, exists := model.SessionList[cookie]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "session not found",
		})
		return
	}

	codeVerifier := c.PostForm("code_verifier")
	expectedChallenge := base64URLEncode(codeVerifier)

	if session.CodeChallenge != expectedChallenge {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid_request. PKCE check failed",
		})
		return
	}

	// アクセストークンの発行
	tokenString := uuid.New().String()
	expireTime := time.Now().Unix() + 3600 // 1時間後に有効期限が切れる

	tokenInfo := model.TokenCode{
		User:      authCode.User,
		ClientId:  authCode.ClientId,
		Scopes:    authCode.Scopes,
		ExpiresAt: expireTime,
	}

	model.TokenCodeList[tokenString] = tokenInfo

	// 認可コードを削除
	delete(model.AuthCodeList, c.PostForm("code"))

	// レスポンスの生成
	tokenResponse := model.TokenResponse{
		AccessToken: tokenString,
		TokenType:   "Bearer",
		ExpiresIn:   expireTime,
	}

	if session.Oidc {
		tokenResponse.IdToken, err = util.MakeJWT(authCode.User)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	resp, err := json.Marshal(tokenResponse)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal error",
		})
		return
	}

	// 正常なレスポンスを返す
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(resp) // JSONオブジェクトを直接レスポンスとして書き込む
}

// Base64 URLエンコード
func base64URLEncode(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
