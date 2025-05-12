package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func DiscordSendTextMessage(webhookURL, message string) {

	payload, err := json.Marshal(map[string]string{
		"content": message,
	})
	if err != nil {
		log.Printf("Discord content marshal error: %v", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Discord send error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.Printf("failed to send log to Discord, status: %s", resp.Status)
	}
}
