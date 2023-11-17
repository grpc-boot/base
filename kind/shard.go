package kind

import "sync"

type Shard[T comparable] interface {
	Set(key T, value any) (oldValue any, exists bool)
	Get(key T) (value any, exists bool)
	Exists(key T) (exists bool)
	Delete(keys ...T) (delNum int)
	Length() int64
}

func NewShard[T comparable](initSize int) Shard[T] {
	return &shard[T]{
		items: make(map[T]any, initSize),
	}
}

type shard[T comparable] struct {
	mutex sync.RWMutex
	items map[T]any
}

func (s *shard[T]) Set(key T, value any) (oldValue any, exists bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldValue, exists = s.items[key]
	s.items[key] = value
	return
}

func (s *shard[T]) Get(key T) (value any, exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists = s.items[key]
	return
}

func (s *shard[T]) Exists(key T) (exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists = s.items[key]
	return
}

func (s *shard[T]) Delete(keys ...T) (delNum int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, key := range keys {
		_, exists := s.items[key]
		if exists {
			delete(s.items, key)
			delNum++
		}
	}

	return
}

func (s *shard[T]) Length() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return int64(len(s.items))
}
