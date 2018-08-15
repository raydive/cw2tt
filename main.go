package cw2tt

import (
	"context"
	"encoding/json"
	"github.com/nulab/go-typetalk/typetalk/v1"
	"golang.org/x/oauth2"
	"net/http"
	"net/url"
	"os"
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
	userAgent := r.Header.Get("User-Agent")
	if !strings.Contains(userAgent, "ChatWork-Webhook") {
		http.Error(w, "You don't have the permission. Only Chatwork Webhook", 401)
	}

	client := makeTypetalkClient()
	client.Messages.PostMessage(context.Background(), 1, "test", nil)
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func makeTypetalkClient() *v1.Client {
	form := url.Values{}
	form.Add("client_id", os.Getenv("CLIENT_ID"))
	form.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	form.Add("grant_type", "client_credentials")
	form.Add("scope", "topic.read,topic.post,topic.write,topic.delete,my")
	oauth2resp, err := http.PostForm("https://typetalk.com/oauth2/access_token", form)
	if err != nil {
		print("Client Credential request returned error")
	}
	v := &AccessToken{}
	json.NewDecoder(oauth2resp.Body).Decode(v)
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: v.AccessToken},
	))
	client := v1.NewClient(tc)
	return client
}
