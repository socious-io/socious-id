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

func organizationsGroup(router *gin.Engine) {
	g := router.Group("organizations")

	g.GET("", paginate(), func(c *gin.Context) {
		paginate, _ := c.MustGet("paginate").(database.Paginate)
		page, limit := c.MustGet("page"), c.MustGet("limit")

		organizations, total, err := models.GetAllOrganizations(paginate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":    page,
			"limit":   limit,
			"results": organizations,
			"total":   total,
		})
	})

	g.GET("/my", paginate(), auth.LoginRequired(), func(c *gin.Context) {
		paginate, _ := c.MustGet("paginate").(database.Paginate)
		user := c.MustGet("user").(*models.User)
		page, limit := c.MustGet("page"), c.MustGet("limit")

		organizations, total, err := models.GetOrganizationsByMember(user.ID, paginate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"page":    page,
			"limit":   limit,
			"results": organizations,
			"total":   total,
		})
	})

	g.GET("/:id", func(c *gin.Context) {
		id := uuid.MustParse(c.Param("id"))

		organization, err := models.GetOrganization(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organization)
	})

	g.POST("", auth.LoginRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		user := c.MustGet("user").(*models.User)

		form := new(OrganizationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization := new(models.Organization)
		utils.Copy(form, organization)

		if err := organization.Create(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := organization.AddMember(ctx, user.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, organization)
	})

	g.PUT("/:id", auth.LoginRequired(), isOrgMember(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		organization := c.MustGet("organization").(*models.Organization)

		form := new(OrganizationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.Copy(form, organization)

		if err := organization.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, organization)
		return
	})

	g.DELETE("/:id", auth.LoginRequired(), isOrgMember(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		organization := c.MustGet("organization").(*models.Organization)

		if err := organization.Remove(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	g.POST("/:id/members/:user_id", auth.LoginRequired(), isOrgMember(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		organization := c.MustGet("organization").(*models.Organization)
		userId := uuid.MustParse(c.Param("user_id"))

		if err := organization.AddMember(ctx, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organization)
		return
	})

	g.DELETE("/:id/members/:user_id", auth.LoginRequired(), isOrgMember(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		organization := c.MustGet("organization").(*models.Organization)
		userId := uuid.MustParse(c.Param("user_id"))

		if err := organization.RemoveMember(ctx, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organization)
		return
	})
}
