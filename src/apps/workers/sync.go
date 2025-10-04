package workers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"socious-id/src/apps/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/socious-io/gomq"
)

func mqSync(_ *models.Access, payload gin.H) {
	gomq.Mq.SendJson("event:identities.sync", payload)
}

func httpSync(access *models.Access, payload gin.H) {
	body, _ := json.Marshal(payload)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, *access.SyncURL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("x-account-center-id", access.ClientID)
	req.Header.Set("x-account-center-secret", access.ClientSecret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

func Sync(userID uuid.UUID) {
	user, err := models.GetUser(userID)
	if err != nil {
		fmt.Println(err)
		return
	}
	organizations, _ := models.GetOrganizationsByMember(user.ID)

	payload := gin.H{
		"user":          user,
		"organizations": organizations,
	}

	accesses, err := models.GetAccesses()
	if err != nil {
		return
	}

	for _, access := range accesses {
		if access.SyncURL == nil {
			continue
		}

		switch access.SyncMethod {
		case models.SyncMethodMq:
			mqSync(&access, payload)
		case models.SyncMethodHttp:
			httpSync(&access, payload)
		}
	}
}
