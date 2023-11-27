package connctx

import "errors"

var (
	ErrType            = errors.New("data format error")
	ErrValueOutOfRange = errors.New("value out of range")
	ErrIndexOutOfRange = errors.New("index out of range")
)
