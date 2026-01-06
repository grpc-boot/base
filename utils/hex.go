package utils

import "encoding/hex"

func HexEncode(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func HexDecode(src []byte) (dst []byte, err error) {
	if len(src)%2 != 0 {
		err = hex.ErrLength
		return
	}

	dst = make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	return dst[:n], err
}

func HexEncode2String(src []byte) string {
	return Bytes2String(HexEncode(src))
}

func HexDecode2String(src []byte) (data string, err error) {
	dst, err := HexDecode(src)
	if err != nil {
		return
	}

	data = Bytes2String(dst)
	return
}
