package kind

import (
	"hash/adler32"
)

var Uint32Hash = adler32.Checksum

// CanHash hash接口
type CanHash interface {
	// HashCode 计算hash值
	HashCode() (hashValue uint32)
}
