package utils

import (
	"encoding/base64"
)

func Base64Encode(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}

func Base64Decode(src []byte) (dst []byte, err error) {
	dst = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	end, err := base64.StdEncoding.Decode(dst, src)
	return dst[:end], err
}

func Base64UrlEncode(src []byte) []byte {
	dst := make([]byte, base64.URLEncoding.EncodedLen(len(src)))
	base64.URLEncoding.Encode(dst, src)
	return dst
}

func Base64UrlDecode(src []byte) (dst []byte, err error) {
	dst = make([]byte, base64.URLEncoding.DecodedLen(len(src)))
	end, err := base64.URLEncoding.Decode(dst, src)
	return dst[:end], err
}

func Base64Encode2String(src []byte) string {
	return Bytes2String(Base64Encode(src))
}

func Base64Decode2String(src []byte) (data string, err error) {
	dst, err := Base64Decode(src)
	if err != nil {
		return
	}

	data = Bytes2String(dst)
	return
}
