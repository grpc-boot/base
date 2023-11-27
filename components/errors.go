package components

import "errors"

var (
	ErrForbidden          = errors.New("permissin denied")
	ErrDataEmpty          = errors.New("data is empty")
	ErrKeyFormat          = errors.New("invalid key format")
	ErrDataFormat         = errors.New("invalid data format")
	ErrDataSign           = errors.New("invalid data sign")
	ErrState              = errors.New("state index must lte 30")
	ErrOutOfRange         = errors.New("out of range")
	ErrAlphanumeric       = errors.New("alphanumeric must be [a-zA-Z0-9] and not repeat")
	ErrAlphanumericLength = errors.New("alphanumeric length must be [50, 62]")
	ErrTimeBack           = errors.New("time go back")
	ErrMachineId          = errors.New("illegal machine id")
)
