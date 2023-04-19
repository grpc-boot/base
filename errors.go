package base

import (
	"google.golang.org/grpc/codes"
)

var (
	ErrForbidden  = NewError(ErrPermissionDenied, "permissin denied")
	ErrDataEmpty  = NewError(ErrDataLoss, "data is empty")
	ErrKeyFormat  = NewError(ErrInvalidArgument, "invalid key format")
	ErrDataFormat = NewError(ErrInvalidArgument, "invalid data format")
	ErrDataSign   = NewError(ErrPermissionDenied, "invalid data sign")
	ErrState      = NewError(ErrOutRange, "state index must lte 30")
)

type MyError struct {
	code codes.Code
	msg  string
}

func NewError(code codes.Code, msg string) error {
	return &MyError{code: code, msg: msg}
}

func (me *MyError) Code() codes.Code {
	return me.code
}

func (me *MyError) Error() string {
	return me.msg
}

func (me *MyError) ToStatus() *Status {
	return StatusWithCode(me.code).WithMsg(me.msg)
}

func (me *MyError) Is(target error) bool {
	tbe, ok := target.(*MyError)
	if !ok {
		return false
	}

	return me.code == tbe.code && me.msg == tbe.msg
}
