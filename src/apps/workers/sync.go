package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"socious-id/src/apps/auth"
	"socious-id/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Sync(userID uuid.UUID) {
	user, err := models.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	organizations, _ := models.GetOrganizationsByMember(user.ID)

	accesses, err := models.GetAccesses()
	if err != nil {
		return
	}

	for _, access := range accesses {
		if access.SyncURL == nil {
			continue
		}
		// FIXME: This is temporary we need to have more control over gorutins
		go syncClient(&access, user, organizations)
	}
}

func syncClient(access *models.Access, user *models.User, organizations []models.Organization) {
	body, _ := json.Marshal(gin.H{
		"user":          user,
		"organizations": organizations,
	})

	clientKey, _ := auth.HashPassword(fmt.Sprintf("%s:%s", access.ClientID, access.ClientSecret))
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, *access.SyncURL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("x-account-center", clientKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
