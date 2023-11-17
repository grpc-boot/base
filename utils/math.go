package utils

import "github.com/grpc-boot/base/v2/kind"

func Abs[V kind.Number](value V) V {
	if value < 0 {
		return -value
	}

	return value
}
