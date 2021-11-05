package base

import (
	"math"
	"sync"

	"github.com/uber-go/atomic"
)

type ShardMap interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) (value interface{}, exists bool)
	Exists(key interface{}) (exists bool)
	Delete(keys ...interface{})
	Length() int64
}

func NewShardMap() ShardMap {
	m := &shardMap{}
	for index := 0; index < math.MaxUint8; index++ {
		m.shardList[index] = shard{
			items: make(map[interface{}]interface{}, 4),
		}
	}

	return m
}

type shardMap struct {
	ShardMap

	shardList [256]shard
	length    atomic.Int64
}

func (m *shardMap) Set(key interface{}, value interface{}) {
	if exists := m.shardList[Index(key)].set(key, value); !exists {
		m.length.Add(1)
	}
}

func (m *shardMap) Get(key interface{}) (value interface{}, exists bool) {
	return m.shardList[Index(key)].get(key)
}

func (m *shardMap) Exists(key interface{}) (exists bool) {
	return m.shardList[Index(key)].exists(key)
}

func (m *shardMap) Delete(keys ...interface{}) {
	for _, key := range keys {
		if exists := m.shardList[Index(key)].delete(key); exists {
			m.length.Sub(1)
		}
	}
}

func (m *shardMap) Length() int64 {
	return m.length.Load()
}

type shard struct {
	mutex sync.RWMutex
	items map[interface{}]interface{}
}

func (s *shard) set(key interface{}, value interface{}) (exists bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists = s.items[key]
	s.items[key] = value
	return
}

func (s *shard) exists(key interface{}) (exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists = s.items[key]
	return
}

func (s *shard) get(key interface{}) (value interface{}, exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists = s.items[key]
	return
}

func (s *shard) delete(key interface{}) (exists bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists = s.items[key]
	if exists {
		delete(s.items, key)
	}
	return
}
