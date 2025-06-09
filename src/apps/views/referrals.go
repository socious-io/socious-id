package views

import (
	"context"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/config"

	"github.com/gin-gonic/gin"
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

		//Search for the referer
		referrerIdentity, err := models.GetReferrerIdentity(form.RefereeID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("couldn't find the referrer user, err: %s\n", err.Error())})
			return
		}
		referralAchievement.ReferrerID = referrerIdentity.ID
		referralAchievement.RewardAmount = getRewardAmountByType(referralAchievement.AchievementType)

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

	g.POST("achievements/claim", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		//Pay the rewards
		//Set all the rewards as claimed
		panic("not implemented")
	})
}

func getRewardAmountByType(t string) float32 {
	rewards := config.Config.ReferralAchievements.Rewards

	for _, reward := range rewards {
		if reward.Type == t {
			return reward.Amount
		}
	}

	return 0
}
