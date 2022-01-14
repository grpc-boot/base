package base

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
)

// HMac HMac算法
func HMac(key []byte, data []byte, hash crypto.Hash) string {
	h := hmac.New(hash.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Hash Hash算法
func Hash(data []byte, hash crypto.Hash) string {
	h := hash.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// Md5 Md5算法
func Md5(data []byte) string {
	return Hash(data, crypto.MD5)
}

// Sha1 Sha1算法
func Sha1(data []byte) string {
	return Hash(data, crypto.SHA1)
}
