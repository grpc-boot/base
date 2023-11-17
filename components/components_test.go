package components

import (
	"github.com/grpc-boot/base/v2/kind"

	"testing"
)

func TestJsonParam_Slice(t *testing.T) {
	jp, _ := UnmarshalJsonParam([]byte(`
		{"cc":[46455, 4545]}
	`))

	cc, _ := jp["cc"].(kind.Slice[int])
	t.Logf("cc:%v", cc)
}
