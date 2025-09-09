// Package notifiers provides functionality for sending notifications to various channels.
package notifiers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kubex-ecosystem/ghbex/internal/defs/interfaces"
)

type Discord struct {
	Webhook string
}

func NewDiscordNotifier(webhook string) interfaces.INotifier {
	return &Discord{
		Webhook: webhook,
	}
}

func (d *Discord) GetType() string {
	return "discord"
}

func (d *Discord) SetWebhook(webhook string) {
	d.Webhook = webhook
}

func (d *Discord) GetWebhook() string {
	return d.Webhook
}

func (d *Discord) Send(ctx context.Context, title, text string, files ...interfaces.IAttachment) error {
	if d.Webhook == "" {
		return nil
	}
	payload := map[string]any{
		"content": "**" + title + "**\n" + text,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", d.Webhook, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	_, err := http.DefaultClient.Do(req)
	return err
}

type Stdout struct{}

func NewStdoutNotifier() interfaces.INotifier {
	return &Stdout{}
}

func (Stdout) GetType() string {
	return "stdout"
}

func (Stdout) SetWebhook(webhook string) {}

func (Stdout) GetWebhook() string {
	return ""
}

func (Stdout) Send(ctx context.Context, title, text string, files ...interfaces.IAttachment) error {
	fmt.Printf("\n==== %s ====\n%s\n", title, text)
	return nil
}
