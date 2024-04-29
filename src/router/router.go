package router

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/NenfuAT/24AuthorizationServer/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init() {
	gin.DisableConsoleColor()
	// ログファイルを作成
	logFile, err := os.Create("log/server.log") // ファイルのパスを指定
	if err != nil {
		fmt.Println("ログファイルの作成に失敗しました:", err)
		return
	}

	// ログの出力先を設定
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout) // ファイルとコンソールにログを出力

	r := gin.Default()
	r.Use(cors.Default())
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!!")
	})

	//認可フロー用
	r.GET("/auth", controller.Auth)
	r.POST("/authcheck", controller.AuthCheck)
	r.POST("/token", controller.TokenHandler)
	r.GET("/userinfo", controller.UserinfoHandler)
	r.GET("/certs", controller.Certs)

	//ユーザ関係
	r.POST("/user/create")

	// サーバの起動とエラーハンドリング
	if err := r.Run("0.0.0.0:8084"); err != nil {
		fmt.Println("サーバの起動に失敗しました:", err)
	}
}
