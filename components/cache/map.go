package cache

import (
	"time"

	"github.com/tinylib/msgp/msgp"
)

type marshaler interface {
	msgp.Marshaler
	msgp.Unmarshaler
}

type Map map[string]interface{}

func (m Map) Strings(key string) []string {
	value, _ := m[key].([]string)
	return value
}

func (m Map) Bool(key string) bool {
	value, _ := m[key].(bool)
	return value
}

func (m Map) String(key string) string {
	value, _ := m[key].(string)
	return value
}

func (m Map) Slice(key string) []interface{} {
	value, _ := m[key].([]interface{})
	return value
}

func (m Map) Bytes(key string) []byte {
	value, _ := m[key].([]byte)
	return value
}

func (m Map) Int(key string) int64 {
	return Int(m[key])
}

func (m Map) Ints(key string) []int {
	value, _ := m[key].([]int)
	return value
}

func (m Map) Int8s(key string) []int8 {
	value, _ := m[key].([]int8)
	return value
}

func (m Map) Int16s(key string) []int16 {
	value, _ := m[key].([]int16)
	return value
}

func (m Map) Int32s(key string) []int32 {
	value, _ := m[key].([]int32)
	return value
}

func (m Map) Int64s(key string) []int64 {
	value, _ := m[key].([]int64)
	return value
}

func (m Map) Uint(key string) uint64 {
	return Uint(m[key])
}

func (m Map) Uints(key string) []uint {
	value, _ := m[key].([]uint)
	return value
}

func (m Map) Uint8s(key string) []uint8 {
	value, _ := m[key].([]uint8)
	return value
}

func (m Map) Uint16s(key string) []uint16 {
	value, _ := m[key].([]uint16)
	return value
}

func (m Map) Uint32s(key string) []uint32 {
	value, _ := m[key].([]uint32)
	return value
}

func (m Map) Uint64s(key string) []uint64 {
	value, _ := m[key].([]uint64)
	return value
}

func (m Map) Float(key string) float64 {
	return Float(m[key])
}

func (m Map) Float64s(key string) []float64 {
	value, _ := m[key].([]float64)
	return value
}

func (m Map) Float32s(key string) []float32 {
	value, _ := m[key].([]float32)
	return value
}

func (m Map) Time(key string) time.Time {
	value, _ := m[key].(time.Time)
	return value
}

func (m Map) Map(key string) Map {
	value, _ := m[key].(map[string]interface{})
	return Map(value)
}
