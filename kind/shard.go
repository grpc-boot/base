package kind

import "sync"

type Shard[K Key] interface {
	Set(key K, value any) (oldValue any, exists bool)
	Get(key K) (value any, exists bool)
	Exists(key K) (exists bool)
	Delete(keys ...K) (delNum int)
	Length() int64
	Range(handler func(key K, value any) bool) bool
}

func NewShard[K Key](initSize int) Shard[K] {
	return &shard[K]{
		items: make(map[K]any, initSize),
	}
}

type shard[K Key] struct {
	mutex sync.RWMutex
	items map[K]any
}

func (s *shard[K]) Set(key K, value any) (oldValue any, exists bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	oldValue, exists = s.items[key]
	s.items[key] = value
	return
}

func (s *shard[K]) Get(key K) (value any, exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, exists = s.items[key]
	return
}

func (s *shard[K]) Exists(key K) (exists bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists = s.items[key]
	return
}

func (s *shard[K]) Delete(keys ...K) (delNum int) {
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

func (s *shard[K]) Length() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return int64(len(s.items))
}

func (s *shard[K]) Range(handler func(key K, value any) bool) bool {
	items := make(map[K]any, len(s.items))
	s.mutex.RLock()
	if len(s.items) == 0 {
		s.mutex.RUnlock()
		return true
	}

	for k, v := range s.items {
		items[k] = v
	}
	s.mutex.RUnlock()

	for key, value := range items {
		if !handler(key, value) {
			items = nil
			return false
		}
	}

	items = nil
	return true
}
