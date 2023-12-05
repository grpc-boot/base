package convert

import (
	"reflect"

	"github.com/grpc-boot/base/v2/internal"
)

func parseType(t reflect.Type, tagName string) []cacheType {
	var (
		length = t.NumField()
		bucket = make([]cacheType, length, length)
		tag    string
	)

	for index := 0; index < length; index++ {
		field := t.Field(index)

		if tagName == "" {
			tag = internal.LcFirst(field.Name)
		} else {
			tag = field.Tag.Get(tagName)
			if tag == "" {
				tag = internal.LcFirst(field.Name)
			}
		}

		bucket[index] = cacheType{
			kind: field.Type.Kind(),
			name: tag,
		}
	}

	return bucket
}
func slice2Map(bucket []cacheType) map[string]cacheMapType {
	data := make(map[string]cacheMapType, len(bucket))
	if len(bucket) < 1 {
		return data
	}

	for index, ct := range bucket {
		data[ct.name] = cacheMapType{
			cacheType: ct,
			index:     index,
		}
	}

	return data
}
