package messaging

type Encoder func(Message) ([]byte, error)

type Decoder func([]byte) (any, error)

