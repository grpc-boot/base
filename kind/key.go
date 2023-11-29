package kind

import (
	"math"
	"reflect"

	"github.com/grpc-boot/base/v2/internal"
)

type Key interface {
	~string | Integer
}

func KeyHash[T Key](t T) uint32 {
	value := reflect.ValueOf(t)
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint32(value.Int() & math.MaxUint32)
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return uint32(value.Uint() & math.MaxUint32)
	default:
		return Uint32Hash(internal.String2Bytes(value.String()))
	}
}
