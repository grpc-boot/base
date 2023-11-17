package utils

import (
	"crypto"
	"crypto/hmac"
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
