package notifier

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func SendSlackAlert(webhookURL, message string) error {
	payload := SlackMessage{Text: message}
	body, _ := json.Marshal(payload)

	_, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(body))
	return err
}
