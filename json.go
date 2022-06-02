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
func JsonEncode(v interface{}) (data []byte, err error) {
	return JsonMarshal(v)
}

// JsonDecode ---
func JsonDecode(data []byte, v interface{}) (err error) {
	return JsonUnmarshal(data, v)
}

// JsonDecodeFile ---
func JsonDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return JsonUnmarshal(conf, v)
}
