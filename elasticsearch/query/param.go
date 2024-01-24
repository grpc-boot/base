package query

import "github.com/grpc-boot/base/v2/kind"

type Param interface {
	Param() kind.JsonParam
}
