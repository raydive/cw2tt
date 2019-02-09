package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/raydive/cw2tt/chatwork"
	"github.com/raydive/cw2tt/typetalk"
)

func main() {
	http.HandleFunc("/send", sendToTypetalk)
	http.HandleFunc("/account", accountHandler)
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

func sendToTypetalk(w http.ResponseWriter, r *http.Request) {
	userAgent := r.UserAgent()
	if !strings.Contains(userAgent, "ChatWork-Webhook") {
		http.Error(w, "You don't have the permission. Only Chatwork Webhook", http.StatusUnauthorized)
		log.Printf("You don't have the permission. Only Chatwork Webhook")
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body.", http.StatusBadRequest)
		log.Print(err)
		return
	}

	// Verify signature
	signature := r.Header.Get("X-ChatWorkWebhookSignature")
	if signature == "" {
		http.Error(w, "You don't have X-ChatWorkWebhookSignature.", http.StatusUnauthorized)
		log.Printf("You don't have X-ChatWorkWebhookSignature.")
		return
	}

	token := os.Getenv("WEBHOOK_TOKEN")
	digest, err := chatwork.GetDigest(token, body)
	if err != nil {
		http.Error(w, "Webhook token is missing. You should check server's environment variables.", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if digest != signature {
		http.Error(w, "Invalid signature.", http.StatusUnauthorized)
		log.Printf("Invalid signature")
		return
	}

	// Take JSON data
	data := chatwork.MessageCreatedWebhook{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Request's body is not excepted JSON.", http.StatusBadRequest)
		log.Print(err)
		return
	}

	topicID, err := strconv.Atoi(os.Getenv("TOPIC_ID"))
	if err != nil {
		http.Error(w, "Invalid topic id. You should check server's environment variables.", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if !data.HasZoomusURI() {
		log.Printf("This post doesn't have zoom.us url.")
		return
	}

	typetalkToken := os.Getenv("TYPETALK_TOKEN")
	ttClient := typetalk.MakeTypetalkBot(typetalkToken)
	_, _, err = ttClient.Messages.PostMessage(context.Background(), topicID, data.Message(), nil)
	if err != nil {
		http.Error(w, "We could not post to Typetalk. You should check server.", http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Invalid body from Typetalk.", http.StatusBadRequest)
		log.Print(err)
		return
	}

	if r.Method == http.MethodGet {

	} else if r.Method == http.MethodPost {
		// TODO must be adding or deleting account name on DB
	}
}
