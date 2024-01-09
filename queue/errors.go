package queue

import "errors"

var (
	ErrNoneHandler = errors.New(`not found handler, please register handler`)
)
