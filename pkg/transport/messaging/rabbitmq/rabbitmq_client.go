package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rasteiro11/PogCore/pkg/config"
)

func newProvider(ctx context.Context) (*amqp.Connection, error) {
	url := config.Instance().RequiredString("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
