package cache

import (
	"os"
	"sync"
)

type bucket struct {
	mutex sync.RWMutex
	items map[string]*Item
}

func (b *bucket) set(key string, item *Item) (isCreate bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.items == nil {
		b.items = map[string]*Item{
			key: item,
		}
		return true
	}

	_, exists := b.items[key]
	b.items[key] = item

	return !exists
}

func (b *bucket) get(key string) (item *Item, exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.items == nil {
		return
	}

	item, exists = b.items[key]
	return
}

func (b *bucket) exists(key string) (exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.items == nil {
		return
	}

	_, exists = b.items[key]
	return
}

func (b *bucket) delete(keys ...string) (delNum int64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.items == nil {
		return
	}

	for _, key := range keys {
		if _, exists := b.items[key]; exists {
			delNum++
		}
		delete(b.items, key)
	}

	return
}

func (b *bucket) flushFile(fileName string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if len(b.items) < 1 {
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	bkt := Bucket{Items: b.items}
	data := make([]byte, 0, bkt.Msgsize())
	data, err = bkt.MarshalMsg(data)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	return err
}
