package convert

import "reflect"

type cacheType struct {
	kind reflect.Kind
	name string
}

type cacheMapType struct {
	cacheType
	index int
}
