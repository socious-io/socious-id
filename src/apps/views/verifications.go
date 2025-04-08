package views

import (
	"context"
	"net/http"
	"net/url"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func verificationsGroup(router *gin.Engine) {
	g := router.Group("verifications")

	g.GET("", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		u, _ := c.MustGet("user").(*models.User)

		v, err := models.GetVerificationByUser(u.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		v.ProofVerify(ctx)

		if v.Status == models.VerificationStatusVerified {
			if err := u.Verify(ctx, models.UserVerificationTypeIdenity); err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "user is verified but couldn't verify user",
				})
				return
			}
		}

		c.JSON(http.StatusOK, v)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		u, _ := c.MustGet("user").(*models.User)

		v := new(models.VerificationCredential)
		v.UserID = u.ID

		if err := v.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.GET("/:id/connect", func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		id := uuid.MustParse(c.Param("id"))

		v, err := models.GetVerification(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if v.ConnectionURL != nil {
			if time.Since(*v.ConnectionAt) < 2*time.Minute {
				c.JSON(http.StatusOK, v)
				return
			}
		}

		callback, _ := url.JoinPath(config.Config.Host, strings.ReplaceAll(c.Request.URL.String(), "connect", "callback"))

		if err := v.NewConnection(ctx, callback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, v)
	})

	g.GET("/:id/callback", func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		id := uuid.MustParse(c.Param("id"))

		v, err := models.GetVerification(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := v.ProofRequest(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
