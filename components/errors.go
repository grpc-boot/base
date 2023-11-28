package components

import "errors"

var (
	ErrDataEmpty          = errors.New("data is empty")
	ErrDataFormat         = errors.New("invalid data format")
	ErrOutOfRange         = errors.New("out of range")
	ErrAlphanumeric       = errors.New("alphanumeric must be [a-zA-Z0-9] and not repeat")
	ErrAlphanumericLength = errors.New("alphanumeric length must be [50, 62]")
	ErrTimeBack           = errors.New("time go back")
	ErrMachineId          = errors.New("illegal machine id")
)
