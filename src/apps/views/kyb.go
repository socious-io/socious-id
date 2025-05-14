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

		orgMembers, err := models.GetOrganizationMembers(org.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, m := range orgMembers {
			//TODO: use nats
			go workers.Sync(m.UserID)
		}

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

		if err := org.Update(ctx); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		orgMembers, err := models.GetOrganizationMembers(org.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for _, m := range orgMembers {
			//TODO: use nats
			go workers.Sync(m.UserID)
		}

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

}

func createDiscordReviewMessage(kyb *models.KYBVerification, u *models.User, org *models.Organization) string {

	documents := ""
	for i, document := range kyb.Documents {
		documents = fmt.Sprintf("%s\n%v. %s", documents, i, document.Url)
	}

	message := fmt.Sprintf("ID: %s\n", kyb.ID)
	message += "\nUser--------------------------------\n"
	message += fmt.Sprintf("ID: %s\n", u.ID)

	if u.FirstName != nil {
		message += fmt.Sprintf("Firstname: %s\n", *u.FirstName)
	} else {
		message += "Firstname: N/A\n"
	}

	if u.LastName != nil {
		message += fmt.Sprintf("Lastname: %s\n", *u.LastName)
	} else {
		message += "Lastname: N/A\n"
	}

	message += fmt.Sprintf("Email: %s\n", u.Email)
	message += "\nOrganization------------------------\n"
	message += fmt.Sprintf("ID: %s\n", org.ID)
	message += fmt.Sprintf("Name: %v\n", org.Name)
	message += fmt.Sprintf("Description: %v\n", org.Description)
	message += fmt.Sprintf("\nDocuments---------------------------%s\n\n", documents)
	message += "\nReviewing----------------------------\n"
	message += fmt.Sprintf("Approve: <%s/kybs/%s/approve?admin_access_token=%s>\n", config.Config.Host, kyb.ID, config.Config.AdminToken)
	message += fmt.Sprintf("Reject: <%s/kybs/%s/reject?admin_access_token=%s>\n", config.Config.Host, kyb.ID, config.Config.AdminToken)

	return message

}
