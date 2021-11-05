package base

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func YamlEncode(v interface{}) (data []byte, err error) {
	return yaml.Marshal(v)
}

func YamlDecode(data []byte, v interface{}) (err error) {
	return yaml.Unmarshal(data, v)
}

func YamlDecodeFile(filePath string, v interface{}) (err error) {
	conf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(conf, v)
}
