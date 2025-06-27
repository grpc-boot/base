package components

import (
	"errors"
	"os"
	"strings"

	"github.com/grpc-boot/base/v3/utils"
)

const (
	ZhCn = `zh_CN`
	En   = `en`
)

var (
	DefaultLang = ZhCn
)

var (
	ErrMsgMapEmpty = errors.New("msgMap can not empty")
	ErrNoYamlFiles = errors.New("there are no `.yml` or `.yaml` files in this directory")
	ErrNoJsonFiles = errors.New("there are no `.json` files in this directory")
)

// I18n 多语言接口
type I18n interface {
	T(key string) (msg string)
	Tl(key, language string) (msg string)
}

// NewI18n 实例化多语言
func NewI18n(msgMap map[string]map[string]string) (i I18n, err error) {
	if len(msgMap) < 1 {
		return nil, ErrMsgMapEmpty
	}

	i = &i18n{
		msgMap: msgMap,
	}

	return i, nil
}

// NewI18nFromYaml 从Yaml文件配置实例化多语言
func NewI18nFromYaml(dir string) (i I18n, err error) {
	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(fileList) < 1 {
		return nil, ErrNoYamlFiles
	}

	var msgMap = make(map[string]map[string]string, len(fileList))
	for _, fi := range fileList {
		fileName := strings.SplitN(fi.Name(), ".", 2)
		if len(fileName) != 2 {
			continue
		}

		if fileName[1] != "yml" && fileName[1] != "yaml" {
			continue
		}

		var msg map[string]string
		err = utils.YamlUnmarshalFile(dir+fi.Name(), &msg)
		if err != nil {
			return nil, err
		}

		msgMap[fileName[0]] = msg
	}

	if len(msgMap) < 1 {
		return nil, ErrNoYamlFiles
	}

	return NewI18n(msgMap)
}

// NewI18nFromJson 从Json文件配置实例化多语言
func NewI18nFromJson(dir string) (i I18n, err error) {
	fileList, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	if len(fileList) < 1 {
		return nil, ErrNoJsonFiles
	}

	var msgMap = make(map[string]map[string]string, len(fileList))
	for _, fi := range fileList {
		fileName := strings.SplitN(fi.Name(), ".", 2)
		if len(fileName) != 2 {
			continue
		}

		if fileName[1] != "json" {
			continue
		}

		var msg map[string]string
		err = utils.JsonUnmarshalFile(dir+fi.Name(), &msg)
		if err != nil {
			return nil, err
		}

		msgMap[fileName[0]] = msg
	}

	if len(msgMap) < 1 {
		return nil, ErrNoJsonFiles
	}
	return NewI18n(msgMap)
}

type i18n struct {
	I18n

	msgMap map[string]map[string]string
}

func (i *i18n) T(key string) string {
	return i.Tl(key, DefaultLang)
}

func (i *i18n) Tl(key, language string) string {
	if _, exists := i.msgMap[language]; !exists {
		return key
	}

	if msg, ok := i.msgMap[language][key]; ok {
		return msg
	}

	return key
}
