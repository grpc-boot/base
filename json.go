package base

import (
	"io/ioutil"

	jsoniter "github.com/json-iterator/go"
)

func JsonEncode(v interface{}) (data []byte, err error) {
	return jsoniter.Marshal(v)
}

func JsonDecode(data []byte, v interface{}) (err error) {
	return jsoniter.Unmarshal(data, v)
}

func JsonDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return jsoniter.Unmarshal(conf, v)
}
