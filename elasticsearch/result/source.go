package result

import (
	"github.com/grpc-boot/base/v3/kind"
	"github.com/grpc-boot/base/v3/utils"
)

type Source kind.JsonParam

func (s Source) Marshal() []byte {
	data, _ := utils.JsonMarshal(s)
	return data
}
