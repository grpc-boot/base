package basis

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var (
	_cache sync.Map
)

type Mapping struct {
	fields   []string
	names    []string
	defaults []string
	labels   []string
	enums    [][]string
	sizes    []int
	kinds    []reflect.Kind
}

func (om *Mapping) Names() []string {
	return om.names
}

func (om *Mapping) Field(name string) string {
	for index, n := range om.names {
		if n == name {
			return om.fields[index]
		}
	}

	return ""
}

func (om *Mapping) Kind(name string) reflect.Kind {
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
	for index, f := range om.fields {
		if f == field {
			return om.fields[index]
		}
	}

	return ""
}

func (om *Mapping) Label(name string) string {
	for index, n := range om.names {
		if n == name {
			return om.labels[index]
		}
	}

	return ""
}

func (om *Mapping) Default(name string) string {
	for index, n := range om.names {
		if n == name {
			return om.defaults[index]
		}
	}

	return ""
}

func (om *Mapping) Size(name string) int {
	for index, n := range om.names {
		if n == name {
			return om.sizes[index]
		}
	}

	return 0
}

func (om *Mapping) Enums(name string) []string {
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
		fieldLen = tp.NumField()
		value    = &Mapping{
			fields:   make([]string, fieldLen),
			names:    make([]string, fieldLen),
			defaults: make([]string, fieldLen),
			labels:   make([]string, fieldLen),
			enums:    make([][]string, fieldLen),
			sizes:    make([]int, fieldLen),
			kinds:    make([]reflect.Kind, fieldLen),
		}
	)

	for i := 0; i < fieldLen; i++ {
		field := tp.Field(i)
		value.names[i] = field.Name
		value.kinds[i] = field.Type.Kind()
		value.fields[i] = field.Tag.Get("field")
		value.labels[i] = field.Tag.Get("label")
		value.defaults[i] = field.Tag.Get("default")
		value.sizes[i], _ = strconv.Atoi(field.Tag.Get("size"))

		enums := field.Tag.Get("enums")
		if len(enums) > 0 {
			value.enums[i] = strings.Split(enums, ",")
		} else {
			value.enums[i] = nil
		}
	}

	_cache.Store(key, value)
	return value
}
