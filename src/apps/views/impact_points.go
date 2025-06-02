package views

import (
	"context"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/apps/workers"

	"github.com/gin-gonic/gin"
	database "github.com/socious-io/pkg_database"
)

func impactPointsGroup(router *gin.Engine) {
	g := router.Group("impact-points")

	g.POST("", clientSecretRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		form := new(ImpactPointForm)
		if err := c.ShouldBind(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		impactPoint := new(models.ImpactPoint)
		utils.Copy(form, impactPoint)

		if err := impactPoint.Create(ctx); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// FIXME: use nats
		go workers.Sync(impactPoint.UserID)

		c.JSON(http.StatusCreated, impactPoint)
	})

	g.GET("", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		impactPoints, total, err := models.GetImpactPoints(user.ID, c.MustGet("paginate").(database.Paginate))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"impact_points": impactPoints,
			"total_count":   total,
			"page":          c.MustGet("page"),
			"limit":         c.MustGet("limit"),
		})
	})

	g.GET("/overview", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		countsPerType, err := models.GetImpactPointsCountsPerType(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"overviews": countsPerType})
	})

	g.GET("badges", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		badges, err := models.GetImpactBadges(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"badges": badges})
	})
}
