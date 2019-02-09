package chatwork

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
)

var ZOOMUS_URL string

func init() {
	ZOOMUS_URL = os.Getenv("ZOOMUS_URL")
}

type MessageCreatedWebhook struct {
	WebhookSettingID string `json:"webhook_setting_id"`
	WebhookEventType string `json:"webhook_event_type"`
	WebhookEventTime int    `json:"webhook_event_time"`
	WebhookEvent     struct {
		MessageID  string `json:"message_id"`
		RoomID     int    `json:"room_id"`
		AccountID  int    `json:"account_id"`
		Body       string `json:"body"`
		SendTime   int    `json:"send_time"`
		UpdateTime int    `json:"update_time"`
	} `json:"webhook_event"`
}

func (hook *MessageCreatedWebhook) HasZoomusURI() bool {
	return strings.Contains(hook.WebhookEvent.Body, ZOOMUS_URL)
}

func (hook *MessageCreatedWebhook) Message() string {
	return "From Chatwork:\n" + strings.Replace(hook.WebhookEvent.Body, "[toall]", "@here", 1)
}

func GetDigest(token string, body []byte) (string, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, decodedToken)
	mac.Write(body)
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return digest, nil
}
