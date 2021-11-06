package base

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
)

func HMac(key []byte, data []byte, hash crypto.Hash) string {
	h := hmac.New(hash.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Hash(data []byte, hash crypto.Hash) string {
	h := hash.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5(data []byte) string {
	return Hash(data, crypto.MD5)
}

func Sha1(data []byte) string {
	return Hash(data, crypto.SHA1)
}
