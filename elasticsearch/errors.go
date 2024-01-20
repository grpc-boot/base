package elasticsearch

import "errors"

var (
	ErrArgsEmpty       = errors.New("args is empty")
	ErrPropertiesEmpty = errors.New("properties is empty")
)
