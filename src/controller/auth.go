package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/NenfuAT/24AuthorizationServer/model"

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
