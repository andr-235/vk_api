package messages

import (
	"context"
	"errors"
	"fmt"

	"github.com/andr-235/vk_api/pkg/client"
)

type MessagesSendParams struct {
	UserID   int    `url:"user_id,omitempty"`
	PeerID   int    `url:"peer_id,omitempty"`
	Domain   string `url:"domain,omitempty"`
	ChatID   int    `url:"chat_id,omitempty"`
	RandomID int    `url:"random_id"`
	Message  string `url:"message,omitempty"`
}

// Validate проверяет валидность параметров метода Send
func (p MessagesSendParams) Validate() error {
	if p.RandomID == 0 {
		return errors.New("random_id обязателен и не должен быть 0")
	}
	if p.UserID == 0 && p.PeerID == 0 && p.ChatID == 0 && p.Domain == "" {
		return errors.New("один из параметров user_id, peer_id, chat_id или domain обязателен")
	}
	return nil
}

func Send(ctx context.Context, c client.Caller, params MessagesSendParams) (int, error) {
	if err := params.Validate(); err != nil {
		return 0, fmt.Errorf("messages.Send: invalid params: %w", err)
	}
	var messageID int
	if err := c.Call(ctx, "messages.send", params, &messageID); err != nil {
		return 0, err
	}
	return messageID, nil
}
