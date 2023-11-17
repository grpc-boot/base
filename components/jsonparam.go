package components

import (
	jsoniter "github.com/json-iterator/go"
)

type JsonParam map[string]any

// UnmarshalJsonParam _
func UnmarshalJsonParam(data []byte) (JsonParam, error) {
	var jp map[string]any
	err := jsoniter.Unmarshal(data, &jp)
	return jp, err
}

// Exists key是否存在，返回bool值
func (jp JsonParam) Exists(key string) bool {
	_, exists := jp[key]
	return exists
}

// String 获取string值，如果不是string类型，返回""
func (jp JsonParam) String(key string) string {
	value, _ := jp[key].(string)
	return value
}

// StringSlice 获取[]string，如果不是[]string返回nil
func (jp JsonParam) StringSlice(key string) []string {
	value, ok := jp[key].([]interface{})
	if !ok {
		return nil
	}

	var stringList = make([]string, 0, len(value))
	if len(value) < 1 {
		return stringList
	}

	for index := 0; index < len(value); index++ {
		if val, ok := value[index].(string); ok {
			stringList = append(stringList, val)
		}
	}

	return stringList
}

// Int 获取Int，如果不是int，返回0
func (jp JsonParam) Int(key string) int {
	return int(jp.Int64(key))
}

// IntSlice 获取[]int，获取失败，返回nil
func (jp JsonParam) IntSlice(key string) []int {
	value, ok := jp[key].([]interface{})
	if !ok {
		return nil
	}

	var intList = make([]int, 0, len(value))
	if len(value) < 1 {
		return intList
	}

	for index := 0; index < len(value); index++ {
		if val, ok := value[index].(float64); ok {
			intList = append(intList, int(val))
		}
	}

	return intList
}

// Int64 获取Int64，如果不是float64，返回0
func (jp JsonParam) Int64(key string) int64 {
	if value, ok := jp[key].(float64); ok {
		return int64(value)
	}

	value, _ := jp[key].(int64)

	return value
}

// Int64Slice 获取[]int64，获取失败返回nil
func (jp JsonParam) Int64Slice(key string) []int64 {
	value, ok := jp[key].([]interface{})
	if !ok {
		return nil
	}

	var int64List = make([]int64, 0, len(value))
	if len(value) < 1 {
		return int64List
	}

	for index := 0; index < len(value); index++ {
		if val, ok := value[index].(float64); ok {
			int64List = append(int64List, int64(val))
		}
	}

	return int64List
}

// Uint32Slice 获取[]uint32，获取失败返回nil
func (jp JsonParam) Uint32Slice(key string) []uint32 {
	value, ok := jp[key].([]interface{})
	if !ok {
		return nil
	}

	var uint32List = make([]uint32, 0, len(value))
	if len(value) < 1 {
		return uint32List
	}

	for index := 0; index < len(value); index++ {
		if val, ok := value[index].(float64); ok {
			uint32List = append(uint32List, uint32(val))
		}
	}

	return uint32List
}

// Float64 获取float64，如果不是float64，返回0
func (jp JsonParam) Float64(key string) float64 {
	value, _ := jp[key].(float64)
	return value
}

// JsonMarshal _
func (jp JsonParam) JsonMarshal() []byte {
	data, _ := jsoniter.Marshal(jp)
	return data
}
