package workers

import (
	"fmt"
	"socious-id/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socious-io/gomq"
)

func Sync(userID uuid.UUID) {
	user, err := models.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	organizations, _ := models.GetOrganizationsByMember(user.ID)

	gomq.Mq.SendJson("event:identities.sync", gin.H{
		"user":          user,
		"organizations": organizations,
	})
}
