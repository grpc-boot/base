package utils

import (
	"crypto"
	"crypto/hmac"
	"fmt"
	"math"

	"github.com/grpc-boot/base/v2/internal"
	"github.com/grpc-boot/base/v2/kind"
)

// HMac HMac算法
func HMac(key []byte, data []byte, hash crypto.Hash) string {
	return Bytes2String(HMacBytes(key, data, hash))
}

func HMacBytes(key []byte, data []byte, hash crypto.Hash) []byte {
	h := hmac.New(hash.New, key)
	h.Write(data)
	return HexEncode(h.Sum(nil))
}

// Hash Hash算法
func Hash(data []byte, hash crypto.Hash) string {
	return Bytes2String(HashBytes(data, hash))
}

func HashBytes(data []byte, hash crypto.Hash) []byte {
	h := hash.New()
	h.Write(data)

	return HexEncode(h.Sum(nil))
}

// Md5 Md5算法
func Md5(data []byte) string {
	return Hash(data, crypto.MD5)
}

// Sha1 Sha1算法
func Sha1(data []byte) string {
	return Hash(data, crypto.SHA1)
}

func Sha256(data []byte) string {
	return Hash(data, crypto.SHA256)
}

// HashValue 计算任意类型hash值
func HashValue(key any) uint32 {
	switch v := key.(type) {
	case kind.CanHash:
		return v.HashCode()
	case bool:
		if v {
			return 1
		}
		return 0
	case uint8:
		return uint32(v)
	case uint16:
		return uint32(v)
	case uint32:
		return key.(uint32)
	case uint64:
		return uint32(v & math.MaxUint32)
	case uint:
		return uint32(v & math.MaxUint32)
	case int8:
		return uint32(v)
	case int16:
		return uint32(v)
	case int32:
		return uint32(v)
	case int64:
		return uint32(v & math.MaxUint32)
	case int:
		return uint32(key.(int) & math.MaxUint32)
	case float64:
		return uint32(int64(v) & math.MaxUint32)
	case float32:
		return uint32(v)
	case string:
		return kind.Uint32Hash(internal.String2Bytes(v))
	case []byte:
		return kind.Uint32Hash(v)
	}

	return kind.Uint32Hash(internal.String2Bytes(fmt.Sprintf("%v", key)))
}

// Index4Bit 索引路由方法，值范围为uint32
func Index4Bit(key any, bitCount uint8) uint32 {
	return HashValue(key) & ((1 << bitCount) - 1)
}

// Index4Uint8 索引路由方法，值范围为uint8
func Index4Uint8(key any) uint8 {
	return uint8(HashValue(key) & math.MaxUint8)
}

// Index4Int8 索引路由方法，值范围为int8
func Index4Int8(key any) int8 {
	return int8(HashValue(key) & math.MaxInt8)
}

// Index4Int16 索引路由方法，值范围为int16
func Index4Int16(key any) int16 {
	return int16(HashValue(key) * math.MaxInt16)
}
