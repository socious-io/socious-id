package views

import (
	"context"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	database "github.com/socious-io/pkg_database"
)

func referralsGroup(router *gin.Engine) {
	g := router.Group("referrals")

	g.GET("", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		paginate := c.MustGet("paginate").(database.Paginate)

		referrals, total, err := models.GetReferrals(identity.ID, paginate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results":     referrals,
			"total_count": total,
			"page":        c.MustGet("page"),
			"limit":       c.MustGet("limit"),
		})
	})

	g.GET("stats", auth.LoginRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		identity := c.MustGet("identity").(*models.Identity)

		referralStats, err := models.GetReferralStats(ctx, identity.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, referralStats)
	})

	g.POST("achievements", clientSecretRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)

		form := new(ReferralAchievementForm)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		referralAchievement := new(models.ReferralAchievement)
		utils.Copy(form, referralAchievement)

		if err := referralAchievement.Create(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, referralAchievement)
	})

	g.GET("achievements", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		identity := c.MustGet("identity").(*models.Identity)
		paginate := c.MustGet("paginate").(database.Paginate)

		referralAchievements, total, err := models.GetReferralAchievements(identity.ID, paginate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"results":     referralAchievements,
			"total_count": total,
			"page":        c.MustGet("page"),
			"limit":       c.MustGet("limit"),
		})
	})

	g.POST("achievements/claim", adminAccessRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		identityId := uuid.MustParse(c.Query("identity_id"))

		_, err := models.GetIdentity(identityId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err = models.ClaimAllReferralAchievements(ctx, identityId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
	})
}
