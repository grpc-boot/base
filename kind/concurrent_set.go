package kind

import "sync"

// concurrentSet hash set，线程安全
type concurrentSet[V comparable] struct {
	data  Set[V]
	mutex sync.RWMutex
}

func NewConcurrentSet[V comparable](initSize uint) Set[V] {
	return &concurrentSet[V]{
		data: NewSet[V](initSize),
	}
}

func (cs *concurrentSet[V]) Add(value V) (isNew bool) {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	return cs.data.Add(value)
}

func (cs *concurrentSet[V]) Del(list ...V) (delNum int) {
	if cs.Length() < 1 {
		return
	}

	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	return cs.data.Del(list...)
}

func (cs *concurrentSet[V]) Exists(value V) (exists bool) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	return cs.data.Exists(value)
}

func (cs *concurrentSet[V]) Length() (length int) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	return cs.data.Length()
}

func (cs *concurrentSet[V]) List() (list Slice[V]) {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()

	return cs.data.List()
}
