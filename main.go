package cw2tt

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"github.com/nulab/go-typetalk/typetalk/v1"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/send", sendToTypetalk)
	http.ListenAndServe(":30000", nil)
}

type webhook struct {
	WebhookSettingID string `json:"webhook_setting_id"`
	WebhookEventType string `json:"webhook_event_type"`
	WebhookEventTime int    `json:"webhook_event_time"`
	WebhookEvent     struct {
		FromAccountID int    `json:"from_account_id"`
		ToAccountID   int    `json:"to_account_id"`
		RoomID        int    `json:"room_id"`
		MessageID     string `json:"message_id"`
		Body          string `json:"body"`
		SendTime      int    `json:"send_time"`
		UpdateTime    int    `json:"update_time"`
	} `json:"webhook_event"`
}

func sendToTypetalk(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	if !strings.Contains(userAgent, "ChatWork-Webhook") {
		http.Error(w, "You don't have the permission. Only Chatwork Webhook", http.StatusUnauthorized)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body.", http.StatusBadRequest)
		return
	}

	// Verify signature
	signature := r.Header.Get("X-ChatWorkWebhookSignature")
	if signature == "" {
		http.Error(w, "You don't have X-ChatWorkWebhookSignature.", http.StatusUnauthorized)
		return
	}

	token := os.Getenv("WEBHOOK_TOKEN")
	digest, err := getDigest(token, body)
	if err != nil {
		http.Error(w, "Webhook token is missing. You should check server's environment variables.", http.StatusInternalServerError)
		return
	}

	if digest != signature {
		http.Error(w, "Invalid signature.", http.StatusUnauthorized)
		return
	}

	// Take JSON data
	data := webhook{}
	err = json.Unmarshal(body, data)
	if err != nil {
		http.Error(w, "Request's body is not excepted JSON.", http.StatusBadRequest)
		return
	}

	topicId, err := strconv.Atoi(os.Getenv("TOPIC_ID"))
	if err != nil {
		http.Error(w, "Invalid topic id. You should check server's environment variables.", http.StatusInternalServerError)
		return
	}

	client := makeTypetalkBot()
	_, _, err = client.Messages.PostMessage(context.Background(), topicId, data.WebhookEvent.Body, nil)
	if err != nil {
		http.Error(w, "We could not post to Typetalk. You should check server.", http.StatusInternalServerError)
		return
	}
}

func getDigest(token string, body []byte) (string, error) {
	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, []byte(decodedToken))
	mac.Write(body)
	digest := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return digest, nil
}

func makeTypetalkBot() *v1.Client {
	token := os.Getenv("TYPETALK_TOKEN")
	client := v1.NewClient(nil)
	client.SetTypetalkToken(token)
	return client
}
