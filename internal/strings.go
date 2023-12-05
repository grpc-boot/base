package internal

import "bytes"

func LcFirst(str string) string {
	strBytes := []byte(str)
	strBytes[0] = bytes.ToLower(strBytes[:1])[0]
	return Bytes2String(strBytes)
}

func UcFirst(str string) string {
	strBytes := []byte(str)
	strBytes[0] = bytes.ToUpper(strBytes[:1])[0]
	return Bytes2String(strBytes)
}
