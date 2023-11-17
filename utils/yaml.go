package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	YamlMarshal   = yaml.Marshal
	YamlUnmarshal = yaml.Unmarshal
)

// YamlEncode ---
func YamlEncode(v any) (data string, err error) {
	bytes, err := YamlMarshal(v)
	if err != nil {
		return "", err
	}

	return Bytes2String(bytes), nil
}

// YamlDecode ---
func YamlDecode(data string, v any) (err error) {
	return YamlUnmarshal(String2Bytes(data), v)
}

// YamlUnmarshalFile ---
func YamlUnmarshalFile(filePath string, v any) (err error) {
	conf, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return YamlUnmarshal(conf, v)
}
