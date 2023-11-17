package kind

import (
	"github.com/grpc-boot/base/v2/internal"
	"hash/crc32"
	"math"

	"go.uber.org/atomic"
)

const shardSize = math.MaxUint8

// ShardMap 分片Map
type ShardMap[T comparable] interface {
	// Set 存储
	Set(key T, value any) (isCreate bool)
	// Get 获取
	Get(key T) (value any, exists bool)
	// Exists 是否存在
	Exists(key T) (exists bool)
	// Delete 删除
	Delete(keys ...T) (delNum int)
	// Length 长度
	Length() int64
	// ShardLength 返回每个分片的长度，可以协助分析元素是否平均分配
	ShardLength() []int64
}

// NewShardMap 实例化ShardMap
func NewShardMap[N Number]() ShardMap[N] {
	m := &shardMap[N]{}
	for index := 0; index <= shardSize; index++ {
		m.shardList[index] = NewShard[N](4)
	}

	return m
}

func NewStringShardMap() *StringShardMap {
	ssm := &StringShardMap{}
	for index := 0; index <= shardSize; index++ {
		ssm.shardList[index] = NewShard[string](4)
	}

	return ssm
}

type shardMap[N Number] struct {
	shardList [shardSize + 1]Shard[N]
	length    atomic.Int64
}

func (m *shardMap[N]) index(key N) int {
	return int(key) & shardSize
}

func (m *shardMap[N]) Set(key N, value any) (isCreate bool) {
	_, exists := m.shardList[m.index(key)].Set(key, value)
	if !exists {
		m.length.Add(1)
	}

	return !exists
}

func (m *shardMap[N]) Get(key N) (value any, exists bool) {
	return m.shardList[m.index(key)].Get(key)
}

func (m *shardMap[N]) Exists(key N) (exists bool) {
	return m.shardList[m.index(key)].Exists(key)
}

func (m *shardMap[N]) Delete(keys ...N) (delNum int) {
	for _, key := range keys {
		if num := m.shardList[m.index(key)].Delete(key); num > 0 {
			m.length.Sub(1)
			delNum++
		}
	}

	return delNum
}

func (m *shardMap[N]) Length() int64 {
	return m.length.Load()
}

func (m *shardMap[N]) ShardLength() []int64 {
	list := make([]int64, len(m.shardList))
	for i, t := range m.shardList {
		list[i] = t.Length()
	}
	return list
}

// StringShardMap string类型shardMap
type StringShardMap struct {
	shardList [shardSize + 1]Shard[string]
	length    atomic.Int64
}

func (sm *StringShardMap) index(key string) int {
	return int(crc32.ChecksumIEEE(internal.String2Bytes(key))) & shardSize
}

func (sm *StringShardMap) Set(key string, value any) (isCreate bool) {
	_, exists := sm.shardList[sm.index(key)].Set(key, value)
	if !exists {
		sm.length.Add(1)
	}

	return !exists
}

func (sm *StringShardMap) Get(key string) (value any, exists bool) {
	return sm.shardList[sm.index(key)].Get(key)
}

func (sm *StringShardMap) Exists(key string) (exists bool) {
	return sm.shardList[sm.index(key)].Exists(key)
}

func (sm *StringShardMap) Delete(keys ...string) (delNum int) {
	for _, key := range keys {
		if num := sm.shardList[sm.index(key)].Delete(key); num > 0 {
			sm.length.Sub(1)
			delNum++
		}
	}

	return delNum
}

func (sm *StringShardMap) Length() int64 {
	return sm.length.Load()
}

func (sm *StringShardMap) ShardLength() []int64 {
	list := make([]int64, len(sm.shardList))
	for i, t := range sm.shardList {
		list[i] = t.Length()
	}
	return list
}
