package utils

import (
	"strings"

	"github.com/grpc-boot/base/v3/kind"
)

func Join[T kind.Integer](sep string, elems ...T) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return StringInteger(elems[0])
	}

	var b strings.Builder
	b.WriteString(StringInteger(elems[0]))
	for _, i := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(StringInteger(i))
	}
	return b.String()
}
