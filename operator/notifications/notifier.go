// Notifier

package notifications

// Notifier sends a notification
type Notifier interface {
	Send(string) error
}
