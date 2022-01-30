package base

import (
	"math"
	"sync"

	"github.com/grpc-boot/base/core/shardmap"

	"go.uber.org/atomic"
)

// ShardMap 分片Map
type ShardMap interface {
	// Set 存储
	Set(key interface{}, value interface{}) (isCreate bool)
	// Get 获取
	Get(key interface{}) (value interface{}, exists bool)
	// Exists 是否存在
	Exists(key interface{}) (exists bool)
	// Delete 删除
	Delete(keys ...interface{}) (delNum int)
	// Length 长度
	Length() int64
}

// NewShardMap 实例化ShardMap
func NewShardMap() ShardMap {
	m := &shardMap{}
	for index := 0; index <= math.MaxUint8; index++ {
		m.shardList[index] = shard{
			items: make(map[interface{}]interface{}, 4),
		}
	}

	return m
}

// NewSharMapWithChannel 实例化ShardMap并返回changeChannel
func NewSharMapWithChannel(size int) (sm ShardMap, ch <-chan shardmap.ChangeEvent) {
	m := &shardMap{}
	for index := 0; index <= math.MaxUint8; index++ {
		m.shardList[index] = shard{
			items: make(map[interface{}]interface{}, 4),
		}
	}

	m.changeChan = make(chan shardmap.ChangeEvent, size)
	return m, m.changeChan
}

type shardMap struct {
	shardList  [256]shard
	length     atomic.Int64
	changeChan chan shardmap.ChangeEvent
}

func (m *shardMap) Set(key interface{}, value interface{}) (isCreate bool) {
	oldValue, exists := m.shardList[Index4Uint8(key)].set(key, value)
	if !exists {
		m.length.Add(1)
	}

	if cap(m.changeChan) < 1 {
		return !exists
	}

	event := shardmap.ChangeEvent{Type: shardmap.Create, Key: key, Value: value}
	if exists {
		event.Type = shardmap.Update
		event.OldValue = oldValue
	}
	m.changeChan <- event
	return !exists
}

func (m *shardMap) Get(key interface{}) (value interface{}, exists bool) {
	return m.shardList[Index4Uint8(key)].get(key)
}

func (m *shardMap) Exists(key interface{}) (exists bool) {
	return m.shardList[Index4Uint8(key)].exists(key)
}

func (m *shardMap) Delete(keys ...interface{}) (delNum int) {
	for _, key := range keys {
		if exists := m.shardList[Index4Uint8(key)].delete(key); exists {
			m.length.Sub(1)
			delNum++

			if cap(m.changeChan) > 0 {
				m.changeChan <- shardmap.ChangeEvent{Type: shardmap.Delete, Key: key}
			}
		}
	}

	return delNum
}

func (m *shardMap) Length() int64 {
	return m.length.Load()
}

type shard struct {
	mutex sync.RWMutex
	items map[interface{}]interface{}
}

func (s *shard) set(key interface{}, value interface{}) (oldValue interface{}, exists bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldValue, exists = s.items[key]
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
