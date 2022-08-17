package base

import (
	"errors"
	"io/ioutil"
	"strings"
)

const (
	ZhCn = `zh_CN`
)

var (
	defaultLan = ZhCn
)

var (
	ErrMsgMapEmpty = errors.New("msgMap can not empty")
	ErrNoYamlFiles = errors.New("there are no `.yml` or `.yaml` files in this directory")
	ErrNoJsonFiles = errors.New("there are no `.json` files in this directory")
)

// I18n 多语言接口
type I18n interface {
	// T 翻译
	T(key, language string) (msg string)
}

// NewI18n 实例化多语言
func NewI18n(defaultLanguage string, msgMap map[string]map[string]string) (i I18n, err error) {
	if defaultLanguage == "" {
		defaultLanguage = defaultLan
	}

	if len(msgMap) < 1 {
		return nil, ErrMsgMapEmpty
	}

	i = &i18n{
		defaultLanguage: defaultLanguage,
		msgMap:          msgMap,
	}

	return i, nil
}

// NewI18nFromYaml 从Yaml文件配置实例化多语言
func NewI18nFromYaml(defaultLanguage string, dir string) (i I18n, err error) {
	fileList, err := ioutil.ReadDir(dir)
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
		err = YamlDecodeFile(dir+fi.Name(), &msg)
		if err != nil {
			return nil, err
		}

		msgMap[fileName[0]] = msg
	}

	if len(msgMap) < 1 {
		return nil, ErrNoYamlFiles
	}

	return NewI18n(defaultLanguage, msgMap)
}

// NewI18nFromJson 从Json文件配置实例化多语言
func NewI18nFromJson(defaultLanguage string, dir string) (i I18n, err error) {
	fileList, err := ioutil.ReadDir(dir)
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
		err = JsonDecodeFile(dir+fi.Name(), &msg)
		if err != nil {
			return nil, err
		}

		msgMap[fileName[0]] = msg
	}

	if len(msgMap) < 1 {
		return nil, ErrNoJsonFiles
	}
	return NewI18n(defaultLanguage, msgMap)
}

type i18n struct {
	I18n

	msgMap          map[string]map[string]string
	defaultLanguage string
}

func (i *i18n) T(key, language string) (msg string) {
	if _, exists := i.msgMap[language]; !exists {
		language = i.defaultLanguage
	}

	msg, _ = i.msgMap[language][key]
	return
}
