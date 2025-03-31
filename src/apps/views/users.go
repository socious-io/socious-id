package views

import (
	"context"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/apps/workers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	})

	g.PUT("/profile", auth.LoginRequired(), func(c *gin.Context) {

		user := c.MustGet("user").(*models.User)
		ctx := c.MustGet("ctx").(context.Context)

		form := new(auth.RegisterForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if _, err := models.GetUserByUsername(form.Username); err == nil {
			c.HTML(http.StatusBadRequest, "update-profile.html", gin.H{
				"error": "Username is already in use. Please select different username.",
			})
			return
		}

		utils.Copy(form, user)

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

		session := sessions.Default(c)

		if session.Get("org_onboard") != nil && session.Get("org_onboard").(bool) {
			c.Redirect(http.StatusSeeOther, "/organizations/register/pre")
			return
		}
		// FIXME: use nats
		go workers.Sync(user.ID)

		c.Redirect(http.StatusSeeOther, "/auth/confirm")
	})

	g.GET("/profile", auth.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "update-profile.html", gin.H{})
	})

	g.PUT("/:id/status", func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)

		form := new(UserUpdateStatusForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := models.GetUser(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := user.UpdateStatus(ctx, form.Status); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	})
}
