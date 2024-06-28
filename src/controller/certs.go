package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kajiLabTeam/mr-platform-authorization-server/util"
)

func Certs(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(util.MakeJWK())
}
