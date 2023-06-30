package query

import "fmt"

var (
	ErrType   = NewErrorf("不支持的数据类型")
	ErrLength = NewErrorf("数据长度错误")
)

type Error struct {
	msg string
}

func (e *Error) Error() string {
	return e.msg
}

func NewErrorf(format string, args ...interface{}) error {
	return &Error{msg: fmt.Sprintf(format, args...)}
}
