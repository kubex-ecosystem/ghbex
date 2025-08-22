// Package common provides common utilities and types.
package common

import (
	"context"
	"fmt"

	"github.com/rafa-mori/ghbex/internal/defs/interfaces"
)

type Notifiers struct {
	Notifiers []*Notifier `mapstructure:",squash"`
}

func NewNotifiersType(nts ...*Notifier) *Notifiers {
	if len(nts) == 0 {
		return &Notifiers{
			Notifiers: []*Notifier{},
		}
	}
	return &Notifiers{
		Notifiers: nts,
	}
}

func NewNotifiers(nts ...interfaces.INotifier) interfaces.INotifiers {
	notifiers := make([]*Notifier, len(nts))
	for i, nt := range nts {
		if n, ok := nt.(*Notifier); ok {
			notifiers[i] = n
		}
	}
	return NewNotifiersType(notifiers...)
}

func (n *Notifiers) GetNotifiers() []interfaces.INotifier {
	if n == nil || len(n.Notifiers) == 0 {
		return nil
	}
	var result []interfaces.INotifier
	for _, notifier := range n.Notifiers {
		result = append(result, notifier)
	}
	return result
}

func (n *Notifiers) AddNotifier(ntf interfaces.INotifier) error {
	if n == nil {
		return nil
	}
	if ntf == nil {
		return fmt.Errorf("notifier cannot be nil")
	}
	if _, ok := ntf.(*Notifier); !ok {
		return fmt.Errorf("notifier must be of type *Notifier")
	}
	n.Notifiers = append(n.Notifiers, ntf.(*Notifier))
	return nil
}

func (n *Notifiers) Send(ctx context.Context, title, text string, files ...interfaces.IAttachment) error {
	for _, notifier := range n.Notifiers {
		if err := notifier.Send(ctx, title, text, files...); err != nil {
			return err
		}
	}
	return nil
}

type Notifier struct {
	Type    string `yaml:"type" json:"type"`
	Webhook string `yaml:"webhook" json:"webhook"`
}

func NewNotifierType(notifierType, webhook string) *Notifier {
	return &Notifier{
		Type:    notifierType,
		Webhook: webhook,
	}
}

func NewNotifier(notifierType, webhook string) interfaces.INotifier {
	return NewNotifierType(notifierType, webhook)
}

func (n *Notifier) GetType() string    { return n.Type }
func (n *Notifier) GetWebhook() string { return n.Webhook }

func (n *Notifier) Send(ctx context.Context, title, text string, files ...interfaces.IAttachment) error {
	// Implementation for sending the notification
	return nil
}
