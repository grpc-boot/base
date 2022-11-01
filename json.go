package base

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

var (
	JsonMarshal   = jsoniter.Marshal
	JsonUnmarshal = jsoniter.Unmarshal
)

// JsonEncode ---
func JsonEncode(v interface{}) (data string, err error) {
	bytes, err := JsonMarshal(v)
	if err != nil {
		return "", err
	}

	return Bytes2String(bytes), nil
}

// JsonDecode ---
func JsonDecode(data string, v interface{}) (err error) {
	return JsonUnmarshal([]byte(data), v)
}

// JsonDecodeFile ---
func JsonDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return JsonUnmarshal(conf, v)
}
