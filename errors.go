package base

import (
	"google.golang.org/grpc/codes"
)

var (
	ErrForbidden          = NewBError(CodePermissionDenied, "permissin denied")
	ErrDataEmpty          = NewBError(CodeDataLoss, "data is empty")
	ErrKeyFormat          = NewBError(CodeInvalidArgument, "invalid key format")
	ErrDataFormat         = NewBError(CodeInvalidArgument, "invalid data format")
	ErrDataSign           = NewBError(CodePermissionDenied, "invalid data sign")
	ErrState              = NewBError(CodeOutRange, "state index must lte 30")
	ErrOutOfRange         = NewBError(CodeOutRange, "out of range")
	ErrAlphanumeric       = NewBError(CodeInvalidArgument, "alphanumeric must be [a-zA-Z0-9] and not repeat")
	ErrAlphanumericLength = NewBError(CodeInvalidArgument, "alphanumeric length must be [50, 62]")
	ErrTimeBack           = NewBError(CodeInternal, "time go back")
	ErrMachineId          = NewBError(CodeInvalidArgument, "illegal machine id")
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
