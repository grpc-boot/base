package components

import (
	"testing"
	"time"

	"github.com/grpc-boot/base/v2/utils"
)

func TestAes_CbcDecrypt(t *testing.T) {
	aes, err := NewAes("b#%*N130js&@1nuc", "3a#%*N130js&@18h")
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	data := time.Now().String() + "撒旦法#$%"
	cryptData := aes.CbcEncrypt([]byte(data))

	t.Logf("base64url: %s", utils.Base64UrlEncode(cryptData))
	t.Logf("hex: %s", utils.HexEncode2String(cryptData))

	explainData, err := aes.CbcDecrypt(cryptData)
	if err != nil {
		t.Fatalf("want nil, got %v", err)
	}

	if utils.Bytes2String(explainData) != data {
		t.Fatalf("want %s, got %s", data, explainData)
	}
}
