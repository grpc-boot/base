package kind

import (
	"go.uber.org/atomic"
	"math"
)

const shardSize = math.MaxUint8

// ShardMap 分片Map
type ShardMap[K Key] interface {
	// Set 存储
	Set(key K, value any) (isCreate bool)
	// Get 获取
	Get(key K) (value any, exists bool)
	// Exists 是否存在
	Exists(key K) (exists bool)
	// Delete 删除
	Delete(keys ...K) (delNum int)
	// Length 长度
	Length() int64
	// ShardLength 返回每个分片的长度，可以协助分析元素是否平均分配
	ShardLength() []int64
	// Range 按照分片遍历元素
	Range(handler func(key K, value any) bool)
}

// NewShardMap 实例化ShardMap
func NewShardMap[K Key]() ShardMap[K] {
	m := &shardMap[K]{}
	for index := 0; index <= shardSize; index++ {
		m.shardList[index] = NewShard[K](4)
	}

	return m
}

type shardMap[K Key] struct {
	shardList [shardSize + 1]Shard[K]
	length    atomic.Int64
}

func (m *shardMap[K]) index(key K) int {
	return int(KeyHash(key) & shardSize)
}

func (m *shardMap[K]) Set(key K, value any) (isCreate bool) {
	_, exists := m.shardList[m.index(key)].Set(key, value)
	if !exists {
		m.length.Add(1)
	}

	return !exists
}

func (m *shardMap[K]) Get(key K) (value any, exists bool) {
	return m.shardList[m.index(key)].Get(key)
}

func (m *shardMap[K]) Exists(key K) (exists bool) {
	return m.shardList[m.index(key)].Exists(key)
}

func (m *shardMap[K]) Delete(keys ...K) (delNum int) {
	for _, key := range keys {
		if num := m.shardList[m.index(key)].Delete(key); num > 0 {
			m.length.Sub(1)
			delNum++
		}
	}

	return delNum
}

func (m *shardMap[K]) Length() int64 {
	return m.length.Load()
}

func (m *shardMap[K]) ShardLength() []int64 {
	list := make([]int64, len(m.shardList))
	for i, t := range m.shardList {
		list[i] = t.Length()
	}
	return list
}

func (m *shardMap[K]) Range(handler func(key K, value any) bool) {
	for _, s := range m.shardList {
		if !s.Range(handler) {
			break
		}
	}
}
