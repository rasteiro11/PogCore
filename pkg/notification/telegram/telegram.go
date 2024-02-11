package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rasteiro11/PogCore/pkg/logger"
)

const providerURL = "https://api.telegram.org/bot%s/sendMessage"

type telegramRequest struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type telegram struct {
	url    *url.URL
	client *http.Client
}

func (t *telegram) SendMessage(ctx context.Context, channel, payload string) ([]byte, error) {
	message, err := json.Marshal(&telegramRequest{
		ChatID:    channel,
		Text:      payload,
		ParseMode: "html",
	})
	if err != nil {
		logger.Global().Errorf("[telegram.SendMessage] json.Marshal() returned error: %+v\n", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx, http.MethodPost, t.url.String(), bytes.NewBuffer(message))
	if err != nil {
		logger.Global().Errorf("[telegram.SendMessage] http.NewRequest() returned error: %+v\n", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		logger.Global().Errorf("[telegram.SendMessage] client.Do() returned error: %+v\n", err)
		return nil, err
	}

	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func NewTelegramClient(token string) (Telegram, error) {
	url, err := url.Parse(fmt.Sprintf(providerURL, token))
	if err != nil {
		logger.Global().Errorf("[telegram.NewTelegramClient] url.Parse() returned error: %+v", err)
		return nil, err
	}

	return &telegram{
		url: url,
		client: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}
