package slack

import "github.com/nlopes/slack"

// Manager manages sending notifications to a Slack channel
type Manager struct {
	client  *slack.Client
	channel string
}

// NewManager creates a new Manager
func NewManager(token string, channel string) *Manager {
	return &Manager{
		client:  slack.New(token),
		channel: channel,
	}
}

// Send sends a message to the slack channel
func (m *Manager) Send(message string) error {
	_, _, err := m.client.PostMessage(
		m.channel,
		slack.MsgOptionText(message, false),
		slack.MsgOptionAsUser(true))
	return err
}
