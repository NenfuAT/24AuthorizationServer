package controller

import (
	"net/http"
	"strings"
	"time"

	"github.com/NenfuAT/24AuthorizationServer/model"
	"github.com/gin-gonic/gin"
)

// ユーザー情報エンドポイント
func UserinfoHandler(c *gin.Context) {
	// Authorizationヘッダーからトークンを取得
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Authorization header missing",
		})
		return
	}

	// Bearerトークンを取得
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) < 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token format",
		})
		return
	}

	accessToken := tokenParts[1]

	// トークンの検証
	tokenInfo, exists := model.TokenCodeList[accessToken]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token not found",
		})
		return
	}

	// トークンの有効期限を確認
	if tokenInfo.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token has expired",
		})
		return
	}

	// スコープの検証
	if tokenInfo.Scopes != "openid profile" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid scope",
		})
		return
	}

	// ユーザー情報を返す
	userInfo := map[string]interface{}{
		"sub":         model.TestUser.ID,
		"email":       model.TestUser.Email,
		"given_name":  model.TestUser.GivenName,
		"family_name": model.TestUser.FamilyName,
		"locale":      model.TestUser.Locale,
	}

	c.JSON(http.StatusOK, userInfo) // ユーザー情報をJSON形式で返す
}
