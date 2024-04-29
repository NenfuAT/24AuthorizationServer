package controller

import (
	"fmt"
	"net/http"

	"github.com/NenfuAT/24AuthorizationServer/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PostUser(c *gin.Context) {
	var req model.User
	if err := c.Bind(&req); err != nil {
		fmt.Println("Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = uuid.New().String()
	user, err := service.CreateUser(req)
	if err != nil {
		fmt.Println("Error creating user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.ID != "" {
		// ユーザーの作成に成功した場合の処理
		c.JSON(http.StatusCreated, user)
	} else {
		// userが空の場合の処理
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email address already in use"})
	}

}
