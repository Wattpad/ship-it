package notifications

import "github.com/nlopes/slack"

// Slack handles notifications to a slack channel
type Slack struct {
	client  *slack.Client
	channel string
}

// NewSlack creates a new slack notifiier
func NewSlack(token string, channel string) *Slack {
	return &Slack{
		client:  slack.New(token),
		channel: channel,
	}
}

// Send sends a message to the slack channel
func (s *Slack) Send(message string) error {
	_, _, err := s.client.PostMessage(
		s.channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true))
	return err
}
