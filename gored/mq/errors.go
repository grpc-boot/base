package mq

import (
	"errors"
	"strings"
)

var (
	ErrGroupEmpty         = errors.New("group cannot be empty")
	ErrConsumerEmpty      = errors.New("consumer cannot be empty")
	ErrConsumerTopicEmpty = errors.New("consumer topic cannot be empty")
)

func IsErrBusyGroup(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), `BUSYGROUP`)
}
