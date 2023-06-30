package base

import (
	"github.com/grpc-boot/base/internal"
)

func JoinInt(sep string, elems ...int) string {
	return internal.JoinInt(sep, elems...)
}

func JoinUint(sep string, elems ...uint) string {
	return internal.JoinUint(sep, elems...)
}

func JoinInt8(sep string, elems ...int8) string {
	return internal.JoinInt8(sep, elems...)
}

func JoinUint8(sep string, elems ...uint8) string {
	return internal.JoinUint8(sep, elems...)
}

func JoinInt16(sep string, elems ...int16) string {
	return internal.JoinInt16(sep, elems...)
}

func JoinUint16(sep string, elems ...uint16) string {
	return internal.JoinUint16(sep, elems...)
}

func JoinInt32(sep string, elems ...int32) string {
	return internal.JoinInt32(sep, elems...)
}

func JoinUint32(sep string, elems ...uint32) string {
	return internal.JoinUint32(sep, elems...)
}

func JoinInt64(sep string, elems ...int64) string {
	return internal.JoinInt64(sep, elems...)
}

func JoinUint64(sep string, elems ...uint64) string {
	return internal.JoinUint64(sep, elems...)
}
