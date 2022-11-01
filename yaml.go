package base

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var (
	YamlMarshal   = yaml.Marshal
	YamlUnmarshal = yaml.Unmarshal
)

// YamlEncode ---
func YamlEncode(v interface{}) (data string, err error) {
	bytes, err := YamlMarshal(v)
	if err != nil {
		return "", err
	}

	return Bytes2String(bytes), nil
}

// YamlDecode ---
func YamlDecode(data string, v interface{}) (err error) {
	return YamlUnmarshal([]byte(data), v)
}

// YamlDecodeFile ---
func YamlDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return YamlUnmarshal(conf, v)
}
