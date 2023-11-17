package kind

import (
	"sync"
)

type Container[K comparable] interface {
	Set(key K, value any)
	Get(key K) (value any, exists bool)
}

func NewContainer[K comparable]() Container[K] {
	return &container[K]{}
}

type container[K comparable] struct {
	data sync.Map
}

func (c *container[K]) Set(key K, value any) {
	c.data.Store(key, value)
}

func (c *container[K]) Get(key K) (value any, exists bool) {
	value, exists = c.data.Load(key)
	return
}
