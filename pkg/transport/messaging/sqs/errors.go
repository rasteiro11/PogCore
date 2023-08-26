package sqs

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

var ErrSizeLimit = errors.New("the size limit to publish is 10 messages")

var ErrNoMessageInQueue = errors.New("no message in queue")

type ErrSendMessage struct {
	Code    string
	Message string
}

var errSendMessage = &ErrSendMessage{}

func (e ErrSendMessage) Error() string {
	return "status=%s message=%s"
}

func (e *ErrSendMessage) Is(target error) bool {
	return errors.Is(target, errSendMessage)
}

func newPublishMessageError(cause types.BatchResultErrorEntry) error {
	return &ErrSendMessage{
		Code:    *cause.Code,
		Message: *cause.Message,
	}
}
