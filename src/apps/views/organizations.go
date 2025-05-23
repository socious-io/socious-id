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
	database "github.com/socious-io/pkg_database"
)

func organizationsGroup(router *gin.Engine) {
	g := router.Group("organizations")

	g.GET("", paginate(), func(c *gin.Context) {
		paginate := c.MustGet("paginate").(database.Paginate)
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

	g.GET("/membered", auth.LoginRequired(), func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)

		organizations, err := models.GetOrganizationsByMember(user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organizations)
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
		// FIXME: use nats
		go workers.Sync(user.ID)
		c.JSON(http.StatusCreated, organization)
	})

	g.PUT("/:id/status", clientSecretRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)

		form := new(OrganizationUpdateStatusForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization, err := models.GetOrganization(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization.Status = form.Status
		if err := organization.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organization)
	})

	g.POST("/:id/verify", clientSecretRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)

		organization, err := models.GetOrganization(uuid.MustParse(c.Param("id")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization.Verified = true
		if err := organization.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organization)
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
		// FIXME: use nats
		go workers.Sync(c.MustGet("user").(*models.User).ID)
		c.JSON(http.StatusAccepted, organization)
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
		// FIXME: use nats
		go workers.Sync(userId)
		c.JSON(http.StatusOK, organization)
	})

	g.DELETE("/:id/members/:user_id", auth.LoginRequired(), isOrgMember(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		organization := c.MustGet("organization").(*models.Organization)
		userId := uuid.MustParse(c.Param("user_id"))

		if userId == c.MustGet("user").(*models.User).ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't remove yourself from the organization"})
			return
		}

		if err := organization.RemoveMember(ctx, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, organization)
	})

	g.GET("/register/pre", auth.LoginRequired(), func(c *gin.Context) {
		next := c.Query("next")
		if next != "" {
			session := sessions.Default(c)
			session.Set("next", next)
			session.Save()
		}
		c.HTML(http.StatusOK, "pre-org-register.html", gin.H{})
	})

	g.GET("/register", auth.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "org-register.html", gin.H{})
	})

	g.POST("/register", auth.LoginRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		user := c.MustGet("user").(*models.User)

		session := sessions.Default(c)

		form := new(OrganizationForm)
		if err := c.ShouldBind(form); err != nil {
			c.HTML(http.StatusBadRequest, "org-register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		organization := new(models.Organization)
		utils.Copy(form, organization)

		if err := organization.Create(ctx); err != nil {
			c.HTML(http.StatusBadRequest, "org-register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := organization.AddMember(ctx, user.ID); err != nil {
			c.HTML(http.StatusBadRequest, "org-register.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		if session.Get("org_onboard") != nil && session.Get("org_onboard").(bool) {
			session.Delete("org_onboard")
			session.Save()
		}

		c.Redirect(http.StatusSeeOther, "/auth/confirm")
	})

	g.GET("/register/complete", auth.LoginRequired(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "post-org-register.html", gin.H{})
	})
}
