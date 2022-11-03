package base

import "errors"

var (
	ErrForbidden  = errors.New("forbidden")
	ErrDataEmpty  = errors.New("data is empty")
	ErrKeyFormat  = errors.New("invalid key format")
	ErrDataFormat = errors.New("invalid data format")
	ErrDataSign   = errors.New("invalid data sign")
)
