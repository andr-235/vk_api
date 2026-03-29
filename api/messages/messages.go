package messages

import (
	"context"

	vk "github.com/andr-235/vk_api"
)

type MessagesSendParams struct {
	UserID   int    `url:"user_id,omitempty"`
	PeerID   int    `url:"peer_id,omitempty"`
	Domain   string `url:"domain,omitempty"`
	ChatID   int    `url:"chat_id,omitempty"`
	RandomID int    `url:"random_id"`
	Message  string `url:"message,omitempty"`
}

func Send(ctx context.Context, c *vk.Client, params MessagesSendParams) (int, error) {
	var messageID int
	if err := c.Call(ctx, "messages.send", params, &messageID); err != nil {
		return 0, err
	}
	return messageID, nil
}
