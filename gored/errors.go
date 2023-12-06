package gored

import "errors"

var (
	ErrInvalidDataType = errors.New("invalid data type")
	ErrNotFound        = errors.New("redis not found")
)
