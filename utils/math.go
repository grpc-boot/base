package utils

import (
	"math"

	"github.com/grpc-boot/base/v2/kind"
)

func Abs[V kind.Number](value V) V {
	if value < 0 {
		return -value
	}

	return value
}

func Ceil[V kind.Number](value float64) V {
	return V(math.Ceil(value))
}
