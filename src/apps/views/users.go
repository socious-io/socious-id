package views

import (
	"context"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"

	"github.com/gin-gonic/gin"
)

func usersGroup(router *gin.Engine) {
	g := router.Group("users")

	g.PUT("/profile", auth.LoginRequired(), func(c *gin.Context) {
		authSession := loadAuthSession(c)
		if authSession == nil {
			c.HTML(http.StatusNotAcceptable, "confirm.html", gin.H{
				"error": "not accepted without auth session",
			})
			return
		}

		u := c.MustGet("user").(*models.User)

		form := new(auth.RegisterForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		u.FirstName = form.FirstName
		u.LastName = form.LastName
		u.Username = *form.Username
		if err := u.UpdateProfile(c.MustGet("ctx").(context.Context)); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		password, _ := auth.HashPassword(*form.Password)
		u.Password = &password
		if err := u.UpdatePassword(c.MustGet("ctx").(context.Context)); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.Redirect(http.StatusSeeOther, "/auth/confirm")
	})

	g.GET("/profile", auth.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "update-profile.html", gin.H{})
	})

}
