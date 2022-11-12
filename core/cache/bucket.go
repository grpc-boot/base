package cache

import (
	"os"
	"sync"
	"time"
)

type bucket struct {
	mutex  sync.RWMutex
	hasNew bool
	data   *Bucket
}

func (b *bucket) set(key string, item *Item) (isCreate bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.hasNew = true

	if b.data == nil {
		b.data = &Bucket{
			Items: map[string]*Item{
				key: item,
			},
		}
		return true
	}

	_, exists := b.data.Items[key]

	b.data.Items[key] = item

	return !exists
}

func (b *bucket) setValue(key string, value []byte) (isCreate bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.hasNew = true

	item := &Item{
		CreatedAt: time.Now().Unix(),
		Value:     value,
	}

	if b.data == nil {
		b.data = &Bucket{
			Items: map[string]*Item{
				key: item,
			},
		}
		return true
	}

	old, exists := b.data.Items[key]
	// 拷贝hit&&miss
	if exists {
		item.Hit = old.Hit
		item.Miss = old.Miss
	}

	b.data.Items[key] = item

	return !exists
}

func (b *bucket) getValue(key string, timeout int64) (value []byte, exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.data == nil {
		return
	}

	var item *Item
	item, exists = b.data.Items[key]
	if !exists {
		return nil, exists
	}

	if item.effective(timeout) {
		return item.Value, exists
	}

	return nil, exists
}

func (b *bucket) get(key string) (item *Item, exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.data == nil {
		return
	}

	item, exists = b.data.Items[key]
	return
}

func (b *bucket) exists(key string) (exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.data == nil {
		return
	}

	_, exists = b.data.Items[key]
	return
}

func (b *bucket) delete(keys ...string) (delNum int64) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.data == nil {
		return
	}

	for _, key := range keys {
		if _, exists := b.data.Items[key]; exists {
			delNum++
			delete(b.data.Items, key)
			b.hasNew = true
		}
	}

	return
}

func (b *bucket) loadFile(fileName string) (loadLength int64, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	if len(data) < 1 {
		return 0, nil
	}

	b.data = &Bucket{Items: map[string]*Item{}}
	_, err = b.data.UnmarshalMsg(data)
	return int64(len(b.data.Items)), err
}

func (b *bucket) flushFile(fileName string) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if !b.hasNew {
		return nil
	}

	if len(b.data.Items) < 1 {
		return nil
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer file.Close()

	data := make([]byte, 0, b.data.Msgsize())
	data, err = b.data.MarshalMsg(data)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err == nil {
		b.hasNew = false
	}
	return err
}
