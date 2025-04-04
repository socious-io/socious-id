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

	g.GET("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, v)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		u, _ := c.Get("user")

		v := new(models.VerificationCredential)
		v.UserID = u.(*models.User).ID

		if err := v.Create(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.GET("/:id/connect", func(c *gin.Context) {
		ctx, _ := c.Get("ctx")
		id := c.Param("id")

		v, err := models.GetVerifications(uuid.MustParse(id))
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

		if err := v.NewConnection(ctx.(context.Context), callback); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, v)
	})

	g.GET("/:id/callback", func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerifications(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := v.ProofRequest(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})

	g.GET("/:id/verify", func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerifications(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, _ := c.Get("ctx")
		v.ProofVerify(ctx.(context.Context))
		c.JSON(http.StatusOK, v)
	})

	g.DELETE("/:id", auth.LoginRequired(), func(c *gin.Context) {
		id := c.Param("id")
		v, err := models.GetVerification(uuid.MustParse(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, _ := c.Get("user")
		if v.UserID.String() != u.(*models.User).ID.String() {
			c.JSON(http.StatusForbidden, gin.H{"error": "not allow"})
			return
		}
		ctx, _ := c.Get("ctx")
		if err := v.Delete(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
