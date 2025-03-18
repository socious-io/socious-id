package views

import (
	"context"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/apps/workers"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")

	g.GET("", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		c.JSON(http.StatusOK, user)
	})

	g.PUT("", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(UserForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		utils.Copy(form, user)

		if err := user.UpdateProfile(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
		return

	})

	g.PUT("/profile", auth.LoginRequired(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		user := c.MustGet("user").(*models.User)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(auth.RegisterForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		user.FirstName = form.FirstName
		user.LastName = form.LastName
		user.AvatarID = form.AvatarID
		if form.Username != nil {
			user.Username = *form.Username
		}
		if err := user.UpdateProfile(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		password, _ := auth.HashPassword(*form.Password)
		user.Password = &password
		if err := user.UpdatePassword(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}
		// FIXME: use nats
		go workers.Sync(user.ID)
		c.Redirect(http.StatusSeeOther, "/auth/confirm")
	})

	g.GET("/profile", auth.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "update-profile.html", gin.H{})
	})

}
