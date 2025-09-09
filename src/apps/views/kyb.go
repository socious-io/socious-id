package views

import (
	"context"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"
	"socious-id/src/apps/utils"
	"socious-id/src/apps/workers"
	"socious-id/src/config"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func kybVerificationGroup(router *gin.Engine) {
	g := router.Group("kybs")

	g.POST("/:id", auth.LoginRequired(), func(c *gin.Context) {
		organizationId := uuid.MustParse(c.Param("id"))
		user := c.MustGet("user").(*models.User)
		ctx := c.MustGet("ctx")

		form := new(KYBVerificationForm)
		if err := c.ShouldBindJSON(form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization, err := models.GetOrganizationByMember(organizationId, user.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Couldn't find organization or you're not a member of"})
			return
		}

		kyb := &models.KYBVerification{
			UserID: user.ID,
			OrgID:  organization.ID,
		}

		if err := kyb.Create(ctx.(context.Context), form.Documents); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		organization.Status = models.OrganizationStatusTypePending
		if err := organization.Update(ctx.(context.Context)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Synchronize
		go workers.Sync(user.ID)

		utils.DiscordSendTextMessage(
			config.Config.Discord.Channel,
			createDiscordReviewMessage(kyb, user, organization),
		)

		c.JSON(http.StatusOK, kyb)
	})

	g.GET("/:id", auth.LoginRequired(), paginate(), func(c *gin.Context) {
		organizationId := uuid.MustParse(c.Param("id"))

		kyb, err := models.GetKybByOrganization(organizationId)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "no kyb found for this identity"})
			return
		}
		c.JSON(http.StatusOK, kyb)
	})

	g.GET("/:id/approve", adminAccessRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		verificationId := c.Param("id")

		verification, err := models.GetKyb(uuid.MustParse(verificationId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := verification.ChangeStatus(ctx, models.KYBStatusApproved); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, _ := models.GetOrganization(verification.OrgID)

		org.Verified = true
		if err := org.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Synchronize
		go workers.Sync(verification.UserID)

		//Add Achievements
		referralAchievement := models.ReferralAchievement{
			RefereeID:       verification.OrgID,
			AchievementType: "KYB",
			Meta: map[string]any{
				"verification": verification,
				"organization": org,
			},
		}
		referralAchievement.Create(ctx)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	g.GET("/:id/reject", adminAccessRequired(), func(c *gin.Context) {
		ctx := c.MustGet("ctx").(context.Context)
		verificationId := c.Param("id")

		verification, err := models.GetKyb(uuid.MustParse(verificationId))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := verification.ChangeStatus(ctx, models.KYBStatusRejected); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		org, _ := models.GetOrganization(verification.OrgID)
		if org.Status != models.OrganizationStatusTypePending {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
			return
		}

		org.Status = models.OrganizationStatusTypeNotActive
		if err := org.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//Synchronize
		go workers.Sync(verification.UserID)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

}

func createDiscordReviewMessage(kyb *models.KYBVerification, u *models.User, org *models.Organization) string {

	documents := ""
	for i, document := range kyb.Documents {
		documents = fmt.Sprintf("%s\n%v. %s", documents, i, document.Url)
	}

	approveUrl := fmt.Sprintf("%s/kybs/%s/approve?admin_access_token=%s", config.Config.Host, kyb.ID, config.Config.AdminToken)
	rejectUrl := fmt.Sprintf("%s/kybs/%s/reject?admin_access_token=%s", config.Config.Host, kyb.ID, config.Config.AdminToken)

	return utils.DedentString(fmt.Sprintf(
		`
			ID: %s

			User--------------------------------
			ID: %s
			Firstname: %s
			Lastname: %s
			Email: %s

			Organization------------------------
			ID: %s
			Name: %v
			Shortname: %s
			Email: %s
			Description: %s

			Documents---------------------------%s

			Reviewing----------------------------
			Approve: <%s>	
			Reject: <%s>
		`,
		kyb.ID,
		u.ID,
		utils.NullableString(u.FirstName),
		utils.NullableString(u.LastName),
		u.Email,
		org.ID,
		utils.NullableString(org.Name),
		org.Shortname,
		utils.NullableString(org.Email),
		utils.NullableString(org.Description),
		documents,
		approveUrl,
		rejectUrl,
	))

}
