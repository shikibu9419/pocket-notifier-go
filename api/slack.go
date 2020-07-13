package api

import (
	"fmt"
	"github.com/slack-go/slack"
	"os"
	"strings"
)

type webhook struct {
	channel string
	message slack.WebhookMessage
}

func NewSlackWebhook(channel string, message slack.WebhookMessage) *webhook {
	return &webhook{
		channel: channel,
		message: message}
}

func (w *webhook) Send() interface{} {
	url := os.Getenv(fmt.Sprintf("SLACK_%s_WEBHOOK_URL", strings.ToUpper(w.channel)))

	return slack.PostWebhook(url, &w.message)
}
