package base

import "sync"

var (
	SetValue = struct{}{}
)

// HashSet 接口
type HashSet interface {
	// Add 添加元素
	Add(items ...interface{}) (newNum int)
	// Delete 删除元素
	Delete(items ...interface{}) (delNum int)
	// Exists 元素是否存在
	Exists(item interface{}) (exists bool)
	// Length 长度
	Length() (length int)
}

// NewHashSet 实例化HashSet
func NewHashSet(size int) HashSet {
	return &hashSet{
		list: make(map[interface{}]struct{}, size),
	}
}

type hashSet struct {
	HashSet

	mutex sync.Mutex
	list  map[interface{}]struct{}
}

func (hs *hashSet) Add(items ...interface{}) (newNum int) {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	for _, item := range items {
		if _, exists := hs.list[item]; !exists {
			hs.list[item] = SetValue
			newNum++
		}
	}
	return
}

func (hs *hashSet) Delete(items ...interface{}) (delNum int) {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	for _, item := range items {
		if _, exists := hs.list[item]; exists {
			delete(hs.list, item)
			delNum++
		}
	}

	return
}

func (hs *hashSet) Exists(item interface{}) (exists bool) {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	_, exists = hs.list[item]
	return
}

func (hs *hashSet) Length() (length int) {
	hs.mutex.Lock()
	defer hs.mutex.Unlock()

	return len(hs.list)
}
