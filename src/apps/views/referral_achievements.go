package views

import (
	"context"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"

	"github.com/gin-gonic/gin"
	database "github.com/socious-io/pkg_database"
)

func referralAchievementsGroup(router *gin.Engine) {
	g := router.Group("referral-achievements")

	g.POST("", clientSecretRequired(), func(c *gin.Context) {
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

		if err := referralAchievement.Create(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, referralAchievement)
	})

	g.GET("", auth.LoginRequired(), paginate(), func(c *gin.Context) {
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
}
