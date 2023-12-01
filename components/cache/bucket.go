package cache

import (
	"os"
	"sync"
	"time"

	"go.uber.org/atomic"
)

type bucket struct {
	mutex  sync.RWMutex
	fMutex sync.Mutex
	hasNew atomic.Bool
	entry  map[string]int64
	data   *Bucket
}

func (b *bucket) setValue(key string, value []byte) (isCreate bool) {
	b.mutex.RLock()
	if b.data != nil {
		if item, exists := b.data.Items[key]; exists {
			item.save(value)

			b.hasNew.Store(true)
			b.mutex.RUnlock()
			return false
		}
	}
	b.mutex.RUnlock()

	b.mutex.Lock()
	defer func() {
		b.hasNew.Store(true)
		b.mutex.Unlock()
	}()

	if b.data == nil {
		b.data = &Bucket{
			Items: map[string]*Item{
				key: {
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
		b.data.Items[key] = &Item{
			CreatedAt:   time.Now().Unix(),
			Value:       value,
			InvokeCount: 1,
		}
	}

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
	defer func() {
		if delNum > 0 {
			b.hasNew.Store(true)
		}
		b.mutex.Unlock()
	}()

	if b.data == nil {
		return
	}

	for _, key := range keys {
		if _, exists := b.data.Items[key]; exists {
			delNum++
			delete(b.data.Items, key)
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
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	if len(data) < 1 {
		return 0, nil
	}

	bkt := &Bucket{Items: map[string]*Item{}}
	_, err = bkt.UnmarshalMsg(data)
	if err != nil {
		return 0, err
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.data = bkt

	return int64(len(b.data.Items)), err
}

func (b *bucket) flushFile(fileName string) (err error) {
	if !b.hasNew.Load() {
		return nil
	}

	b.mutex.RLock()
	if !b.hasNew.Load() || b.data == nil {
		b.mutex.RUnlock()
		return nil
	}

	data := make([]byte, 0, b.data.Msgsize())
	data, err = b.data.MarshalMsg(data)
	if err != nil {
		b.mutex.RUnlock()
		return err
	}
	b.hasNew.Store(false)
	b.mutex.RUnlock()

	b.fMutex.Lock()
	defer b.fMutex.Unlock()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
