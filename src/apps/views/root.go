package views

import (
	"net/http"
	"socious-id/src/apps/shortener"

	"github.com/gin-gonic/gin"
)

func rootGroup(router *gin.Engine) {
	g := router.Group("")

	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	shortener.Routers(g)
}
