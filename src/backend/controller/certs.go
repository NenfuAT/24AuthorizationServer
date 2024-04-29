package controller

import (
	"net/http"

	"github.com/NenfuAT/24AuthorizationServer/util"
	"github.com/gin-gonic/gin"
)

func Certs(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(util.MakeJWK())
}
