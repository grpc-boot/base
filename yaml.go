package base

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// YamlEncode ---
func YamlEncode(v interface{}) (data []byte, err error) {
	return yaml.Marshal(v)
}

// YamlDecode ---
func YamlDecode(data []byte, v interface{}) (err error) {
	return yaml.Unmarshal(data, v)
}

// YamlDecodeFile ---
func YamlDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(conf, v)
}
