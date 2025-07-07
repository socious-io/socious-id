package views

import (
	"context"
	"net/http"
	"net/url"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/workers"
	"socious-id/src/config"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func credentialsGroup(router *gin.Engine) {
	g := router.Group("credentials")

	g.GET("", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		u, _ := c.MustGet("user").(*models.User)

		var credentialType models.CredentialType
		err := credentialType.Scan(c.Query("type"))
		if err != nil {
			credentialType = models.CredentialTypeKYC
		}

		currentVerificationStatus := u.IdentityVerifiedAt

		v, err := models.GetCredentialByUserAndType(u.ID, credentialType)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		v.HandleByType(ctx)

		if v.Type == models.CredentialTypeKYC && v.Status == models.CredentialStatusVerified {
			if err := u.Verify(ctx, models.UserVerificationTypeIdentity); err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{
					"error": "user is verified but couldn't verify user",
				})
				return
			}

			if currentVerificationStatus == nil && u.IdentityVerifiedAt != nil {
				go workers.Sync(u.ID)

				//Add Achievements
				referralAchievement := models.ReferralAchievement{
					RefereeID:       v.UserID,
					AchievementType: "REF_KYC",
					Meta: map[string]any{
						"credential": v,
						"user":       u,
					},
				}
				referralAchievement.Create(ctx)
			}
		}

		c.JSON(http.StatusOK, v)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		u, _ := c.MustGet("user").(*models.User)

		form := new(CredentialForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		v := new(models.Credential)
		v.UserID = u.ID

		if err := v.Create(ctx, form.Type); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, v)
	})

	g.GET("/:id/connect", func(c *gin.Context) {
		ctx, _ := c.MustGet("ctx").(context.Context)
		id := uuid.MustParse(c.Param("id"))

		v, err := models.GetCredential(id)
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

		v, err := models.GetCredential(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := v.HandleByType(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
}
