package telegram

import "context"

type Telegram interface {
	SendMessage(ctx context.Context, channel, payload string) ([]byte, error)
}
