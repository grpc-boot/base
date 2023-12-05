package kind

import (
	"hash/crc32"
)

var Uint32Hash = crc32.ChecksumIEEE

// CanHash hash接口
type CanHash interface {
	// HashCode 计算hash值
	HashCode() (hashValue uint32)
}
