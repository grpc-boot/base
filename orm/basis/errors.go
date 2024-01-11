package basis

import "errors"

var (
	ErrNoSetter   = errors.New("setter is empty")
	ErrDataFormat = errors.New("invalid data format")
)
