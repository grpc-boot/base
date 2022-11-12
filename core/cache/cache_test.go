package cache

import (
	"testing"
)

func TestNew(t *testing.T) {
	cache := New("./")
	cache.Set("key", []byte("value"), 36)
	item, exist := cache.Get("key")
	t.Logf("item:%v, exists:%t", item, exist)
	cache.flushFile()
}
