package notification

import (
	"github.com/slack-go/slack"
)

type SlackNotifier struct {
	webhookURL string
}

func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{webhookURL: webhookURL}
}

func (s *SlackNotifier) Notify(message string) error {
	msg := &slack.WebhookMessage{
		Text: message,
	}
	return slack.PostWebhook(s.webhookURL, msg)
}
