package interfaces

import "context"

type INotifiers interface {
	// AddNotifier adds a new notifier to the collection.
	AddNotifier(ntf INotifier) error
	// GetNotifiers returns a slice of all configured notifiers.
	GetNotifiers() []INotifier
	// Send sends a notification with the given title and text.
	Send(ctx context.Context, title, text string, files ...IAttachment) error
}

// INotifier represents a single notification configuration.
type INotifier interface {
	GetType() string
	GetWebhook() string

	Send(ctx context.Context, title, text string, files ...IAttachment) error
}
