package elasticsearch

import "errors"

var (
	ErrArgsEmpty       = errors.New("args is empty")
	ErrSetterEmpty     = errors.New("setter is empty")
	ErrPropertiesEmpty = errors.New("properties is empty")
	ErrIndexNotExists  = errors.New("no such index")
)
