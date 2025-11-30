package rabbitmq

import "errors"

var ErrNoMessageInQueue = errors.New("no message in queue")
var ErrProviderNotInitialized = errors.New("rabbitmq provider not initialized")
var ErrInvalidDeliveryTag = errors.New("invalid delivery tag")
