package msg

import "time"

type Value interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | bool | string | []byte | time.Time | []interface{} | map[string]string | map[string]interface{} | Map
}
