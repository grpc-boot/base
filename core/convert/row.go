package convert

import (
	"fmt"
	"reflect"
	"strconv"
)

type Row map[string]string

func (r Row) String(key string) string {
	value, _ := r[key]
	return value
}

func (r Row) Float64(key string) float64 {
	value, _ := r[key]
	if value == "" {
		return 0
	}

	val, _ := strconv.ParseFloat(value, 64)
	return val
}

func (r Row) Int64(key string) int64 {
	value, _ := r[key]
	if value == "" {
		return 0
	}

	val, _ := strconv.ParseInt(value, 10, 64)
	return val
}

func (r Row) Uint64(key string) uint64 {
	value, _ := r[key]
	if value == "" {
		return 0
	}

	val, _ := strconv.ParseUint(value, 10, 64)
	return val
}

func (r Row) Int(key string) int {
	return int(r.Int64(key))
}

func (r Row) Uint(key string) uint {
	return uint(r.Uint64(key))
}

func (r Row) Bool(key string) bool {
	value, _ := r[key]
	if value == "" {
		return false
	}

	val, _ := strconv.ParseBool(r[key])
	return val
}

func (r Row) Convert(out interface{}) (err error) {
	return r.ConvertByTag(out, "json")
}

func (r Row) ConvertByTag(out interface{}, tagName string) (err error) {
	var (
		fields      []cacheType
		t           = reflect.TypeOf(out)
		val, exists = _cache.Load(t)
	)

	if !exists {
		if t.Kind() != reflect.Ptr {
			return fmt.Errorf("obj must be a pointer")
		}

		if t.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("obj must be a pointer to a struct")
		}

		fields = parseType(t.Elem(), tagName)
		_cache.Store(t, fields)
		_mapCache.Store(t, slice2Map(fields))
	} else {
		fields, _ = val.([]cacheType)
	}

	value := reflect.ValueOf(out).Elem()
	for index, ct := range fields {
		switch ct.kind {
		case reflect.String:
			value.Field(index).SetString(r.String(ct.name))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value.Field(index).SetInt(r.Int64(ct.name))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			value.Field(index).SetUint(r.Uint64(ct.name))
		case reflect.Float32, reflect.Float64:
			value.Field(index).SetFloat(r.Float64(ct.name))
		case reflect.Bool:
			value.Field(index).SetBool(r.Bool(ct.name))
		default:
			continue
		}
	}

	return nil
}
