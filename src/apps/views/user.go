package views

import (
	"context"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")

	g.GET("/", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		c.JSON(http.StatusOK, user)
	})

	g.PUT("/", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		ctx, _ := c.Get("ctx")

		form := new(UserForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.Copy(form, user)

		if err := user.UpdateProfile(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
		return

	})
}
