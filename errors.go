package base

import (
	"google.golang.org/grpc/codes"
)

var (
	ErrForbidden  = NewError(CodePermissionDenied, "permissin denied")
	ErrDataEmpty  = NewError(CodeDataLoss, "data is empty")
	ErrKeyFormat  = NewError(CodeInvalidArgument, "invalid key format")
	ErrDataFormat = NewError(CodeInvalidArgument, "invalid data format")
	ErrDataSign   = NewError(CodePermissionDenied, "invalid data sign")
	ErrState      = NewError(CodeOutRange, "state index must lte 30")
	ErrOutOfRange = NewError(CodeOutRange, "out of range")
	ErrTimeBack   = NewError(CodeInternal, "time go back")
	ErrMachineId  = NewError(CodeInvalidArgument, "illegal machine id")
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
