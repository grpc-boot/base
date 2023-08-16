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

type BError struct {
	code codes.Code
	msg  string
	flag uint8
}

func NewBError(code codes.Code, msg string) *BError {
	return &BError{code: code, msg: msg}
}

func NewError(code codes.Code, msg string) error {
	return NewBError(code, msg)
}

func (be *BError) Code() codes.Code {
	return be.code
}

func (be *BError) Error() string {
	return be.msg
}

func (be *BError) WithFlag(flag uint8) *BError {
	be.flag = flag
	return be
}

func (be *BError) ToStatus() *Status {
	return StatusWithCode(be.code).WithMsg(be.msg).WithFlag(be.flag)
}

func (be *BError) Is(target error) bool {
	tbe, ok := target.(*BError)
	if !ok {
		return false
	}

	return be.code == tbe.code && be.msg == tbe.msg
}

func (be *BError) Equal(target error) bool {
	tbe, ok := target.(*BError)
	if !ok {
		return false
	}

	return be.code == tbe.code && be.flag == tbe.flag && be.msg == tbe.msg
}

func (be *BError) Clone() *BError {
	return NewBError(be.code, be.msg).WithFlag(be.flag)
}

func (be *BError) JsonMarshal() []byte {
	return be.ToStatus().JsonMarshal()
}
