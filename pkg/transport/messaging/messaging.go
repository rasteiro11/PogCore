package messaging

import "context"

type Attributes interface {
	Add(key, value string)
	Get(key string) string
	Lookup(key string) (string, bool)
	Delete(key string)
	Values() map[string][]string
}

type Message interface {
	Payload() any
	Attributes() Attributes
	ToHash() string
}

type Publisher interface {
	Publish(context.Context, ...Message) error
}

type Subscriber interface {
	Subscribe(context.Context) (Message, error)
	Commit(context.Context, Message) error
	Close(context.Context) error
}
