package basis

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	_cache     sync.Map
	MinMapSize = 128
)

type Mapping struct {
	fields   []string
	names    []string
	defaults []string
	labels   []string
	enums    [][]string
	sizes    []int
	kinds    []reflect.Kind

	fieldsMap   map[string]string
	namesMap    map[string]string
	defaultsMap map[string]string
	labelsMap   map[string]string
	enumsMap    map[string][]string
	sizesMap    map[string]int
	kindsMap    map[string]reflect.Kind
}

func newMapping(size int) *Mapping {
	mp := &Mapping{
		names: make([]string, size),
	}

	if size >= MinMapSize {
		mp.fieldsMap = make(map[string]string, size)
		mp.namesMap = make(map[string]string, size)
		mp.defaultsMap = make(map[string]string, size)
		mp.labelsMap = make(map[string]string, size)
		mp.enumsMap = make(map[string][]string, size)
		mp.sizesMap = make(map[string]int, size)
		mp.kindsMap = make(map[string]reflect.Kind, size)
	} else {
		mp.fields = make([]string, size)
		mp.defaults = make([]string, size)
		mp.labels = make([]string, size)
		mp.enums = make([][]string, size)
		mp.sizes = make([]int, size)
		mp.kinds = make([]reflect.Kind, size)
	}

	return mp
}

func (om *Mapping) Names() []string {
	return om.names
}

func (om *Mapping) Field(name string) string {
	if len(om.fieldsMap) != 0 {
		value, _ := om.fieldsMap[name]
		return value
	}

	for index, n := range om.names {
		if n == name {
			return om.fields[index]
		}
	}

	return ""
}

func (om *Mapping) Kind(name string) reflect.Kind {
	if len(om.fieldsMap) != 0 {
		value, _ := om.kindsMap[name]
		return value
	}

	for index, n := range om.names {
		if n == name {
			return om.kinds[index]
		}
	}

	return reflect.Invalid
}

func (om *Mapping) Index(name string) int {
	for index, n := range om.names {
		if n == name {
			return index
		}
	}

	return -1
}

func (om *Mapping) Name(field string) string {
	if len(om.fieldsMap) != 0 {
		value, _ := om.namesMap[field]
		return value
	}

	for index, f := range om.fields {
		if f == field {
			return om.fields[index]
		}
	}

	return ""
}

func (om *Mapping) Label(name string) string {
	if len(om.fieldsMap) != 0 {
		value, _ := om.labelsMap[name]
		return value
	}

	for index, n := range om.names {
		if n == name {
			return om.labels[index]
		}
	}

	return ""
}

func (om *Mapping) Default(name string) string {
	if len(om.fieldsMap) != 0 {
		value, _ := om.defaultsMap[name]
		return value
	}

	for index, n := range om.names {
		if n == name {
			return om.defaults[index]
		}
	}

	return ""
}

func (om *Mapping) Size(name string) int {
	if len(om.fieldsMap) != 0 {
		value, _ := om.sizesMap[name]
		return value
	}

	for index, n := range om.names {
		if n == name {
			return om.sizes[index]
		}
	}

	return 0
}

func (om *Mapping) Enums(name string) []string {
	if len(om.fieldsMap) != 0 {
		value, _ := om.enumsMap[name]
		return value
	}

	for index, n := range om.names {
		if n == name {
			return om.enums[index]
		}
	}

	return nil
}

func ParseMapping(t any) *Mapping {
	key := ""

	if model, ok := t.(Model); ok {
		// 先根据表名称检索
		key = model.TableName()
		if key != "" {
			om, exists := _cache.Load(key)
			if exists {
				return om.(*Mapping)
			}
		}
	}

	tp := reflect.TypeOf(t)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	if key == "" {
		// 根据结构体检索
		key = fmt.Sprintf("%s.%s", tp.PkgPath(), tp.Name())
		om, exists := _cache.Load(key)
		if exists {
			return om.(*Mapping)
		}
	}

	var (
		fieldLen  = tp.NumField()
		value     = newMapping(fieldLen)
		field     reflect.StructField
		fieldName string
		enums     string
	)

	for i := 0; i < fieldLen; i++ {
		field = tp.Field(i)
		value.names[i] = field.Name

		fieldName = field.Tag.Get("field")
		enums = field.Tag.Get("enums")

		if len(value.fieldsMap) > 0 {
			value.namesMap[fieldName] = field.Name
			value.kindsMap[field.Name] = field.Type.Kind()
			value.fieldsMap[field.Name] = field.Tag.Get("field")
			value.labelsMap[field.Name] = field.Tag.Get("label")
			value.defaultsMap[field.Name] = field.Tag.Get("default")
			value.sizesMap[field.Name], _ = strconv.Atoi(field.Tag.Get("size"))
			if len(enums) > 0 {
				value.enumsMap[field.Name] = strings.Split(enums, ",")
			}
		} else {
			value.kinds[i] = field.Type.Kind()
			value.fields[i] = field.Tag.Get("field")
			value.labels[i] = field.Tag.Get("label")
			value.defaults[i] = field.Tag.Get("default")
			value.sizes[i], _ = strconv.Atoi(field.Tag.Get("size"))
			if len(enums) > 0 {
				value.enums[i] = strings.Split(enums, ",")
			} else {
				value.enums[i] = nil
			}
		}
	}

	_cache.Store(key, value)
	return value
}
