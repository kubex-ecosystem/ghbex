// Package notifiers provides functionality for sending notifications to various channels.
package notifiers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rafa-mori/ghbex/internal/config"
)

type Discord struct {
	Webhook string
}

func (d Discord) Send(ctx context.Context, title, text string, files ...config.Attachment) error {
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

func (Stdout) Send(ctx context.Context, title, text string, files ...config.Attachment) error {
	fmt.Printf("\n==== %s ====\n%s\n", title, text)
	return nil
}
