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
func YamlEncode(v interface{}) (data []byte, err error) {
	return YamlMarshal(v)
}

// YamlDecode ---
func YamlDecode(data []byte, v interface{}) (err error) {
	return YamlUnmarshal(data, v)
}

// YamlDecodeFile ---
func YamlDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return YamlUnmarshal(conf, v)
}
