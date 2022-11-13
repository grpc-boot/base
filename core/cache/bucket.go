package cache

import (
	"os"
	"sync"
	"time"
)

type bucket struct {
	mutex  sync.RWMutex
	hasNew bool
	entry  map[string]int64
	data   *Bucket
}

func (b *bucket) setValue(key string, value []byte) (isCreate bool) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.hasNew = true

	if b.data == nil {
		b.data = &Bucket{
			Items: map[string]*Item{
				key: &Item{
					CreatedAt:   time.Now().Unix(),
					Value:       value,
					InvokeCount: 1,
				},
			},
		}
		return true
	}

	item, exists := b.data.Items[key]

	if exists {
		item.Value = value
		item.CreatedAt = time.Now().Unix()
		item.InvokeCount++
	} else {
		item = &Item{
			CreatedAt:   time.Now().Unix(),
			Value:       value,
			InvokeCount: 1,
		}
	}

	b.data.Items[key] = item

	return !exists
}

func (b *bucket) getValue(key string, timeout int64) (value []byte, lock *int64, effective bool, exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.data == nil {
		return
	}

	var item *Item
	item, exists = b.data.Items[key]
	if !exists {
		return
	}

	effective = item.effective(timeout)

	return item.Value, &item._lock, effective, exists
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

func (b *bucket) items() []Item {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.data == nil || len(b.data.Items) < 1 {
		return nil
	}

	list := make([]Item, 0, len(b.data.Items))
	for key, item := range b.data.Items {
		list = append(list, Item{
			Key:         key,
			Value:       item.Value,
			Hit:         item.Hit,
			Miss:        item.Miss,
			InvokeCount: item.InvokeCount,
			CreatedAt:   item.CreatedAt,
		})
	}

	return list
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

	if b.data == nil {
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
