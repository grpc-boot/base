package utils

import (
	"os"

	"github.com/goccy/go-json"
)

var (
	JsonMarshal   = json.Marshal
	JsonUnmarshal = json.Unmarshal
)

// JsonEncode ---
func JsonEncode(v any) (data string, err error) {
	bytes, err := JsonMarshal(v)
	if err != nil {
		return "", err
	}

	return Bytes2String(bytes), nil
}

// JsonDecode ---
func JsonDecode(data string, v any) (err error) {
	return JsonUnmarshal(String2Bytes(data), v)
}

// JsonUnmarshalFile ---
func JsonUnmarshalFile(filePath string, v any) (err error) {
	conf, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return JsonUnmarshal(conf, v)
}
