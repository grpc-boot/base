package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/tinylib/msgp/msgp"
)

//------------------------ Bucket Start--------------------------------

func (b *Bucket) Marshal() (data []byte, err error) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	data = make([]byte, 0, b.Msgsize())
	return b.MarshalMsg(data)
}

func (b *Bucket) set(key string, value interface{}) (isCreate bool, err error) {
	b.mutex.Lock()
	defer func() {
		b.hasChanged = true
		b.mutex.Unlock()
	}()

	now := time.Now().Unix()
	if b.Data == nil {
		item := &Entry{CreatedAt: now}
		err = item.save(value, now)
		if err != nil {
			return
		}

		isCreate = true
		b.Data = map[string]*Entry{
			key: item,
		}
		return
	}

	entry, exists := b.Data[key]
	if exists {
		err = entry.save(value, now)
		if err != nil {
			return
		}
	} else {
		entry = &Entry{CreatedAt: now}
		err = entry.save(value, now)
		if err != nil {
			return
		}

		isCreate = true
		b.Data[key] = entry
	}
	return
}

func (b *Bucket) get(key string, timeout int64) (value interface{}, lock *int64, effective bool, exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.Data == nil {
		return
	}

	var entry *Entry
	entry, exists = b.Data[key]
	if !exists {
		return
	}

	effective = entry.effective(timeout, time.Now().Unix())
	if effective {
		entry.hitCnt.Add(1)
	} else {
		entry.missCnt.Add(1)
	}

	return entry.Value, &entry._lock, effective, exists
}

func (b *Bucket) exists(key string) (exists bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.Data == nil {
		return
	}

	_, exists = b.Data[key]
	return
}

func (b *Bucket) delete(keys ...string) (delNum int64) {
	b.mutex.Lock()
	defer func() {
		if delNum > 0 {
			b.hasChanged = true
		}
		b.mutex.Unlock()
	}()

	if b.Data == nil {
		return
	}

	for _, key := range keys {
		if _, exists := b.Data[key]; exists {
			delNum++
			delete(b.Data, key)
		}
	}

	return
}

func (b *Bucket) items() []Item {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.Data == nil || len(b.Data) < 1 {
		return nil
	}

	list := make([]Item, 0, len(b.Data))
	for key, entry := range b.Data {
		list = append(list, entry.toItem(key))
	}
	return list
}

func (b *Bucket) loadFile(fileName string) (loadLength int64, err error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, err
	}

	if len(data) < 1 {
		return 0, nil
	}

	bkt := &Bucket{}
	_, err = bkt.UnmarshalMsg(data)
	if err != nil {
		return 0, err
	}

	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.Data = bkt.Data

	return int64(len(b.Data)), err
}

func (b *Bucket) needFlush() bool {
	if time.Now().Unix()-b.latestSync.Load().Unix() > FlushWithoutChangeIntervalSeconds {
		return true
	}

	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.hasChanged {
		return true
	}

	if b.Data == nil {
		return false
	}
	return true
}

func (b *Bucket) flushFile(fileName string) (err error) {
	if !b.needFlush() {
		return nil
	}

	// 允许误差
	b.hasChanged = false

	data, err := b.Marshal()
	if err != nil {
		return err
	}

	// 保证同一个文件没有并发
	b.fMutex.Lock()
	defer b.fMutex.Unlock()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err == nil {
		b.latestSync.Store(time.Now())
	}
	return err
}

//------------------------ Bucket End--------------------------------

//------------------------ Entry Start--------------------------------

func (e *Entry) save(value interface{}, now int64) (err error) {
	if value == nil {
		e.Value = nil
		return
	}

	defer func() {
		if err == nil {
			e.UpdatedAt = now
			e.updateCnt++
		}
	}()

	switch val := value.(type) {
	case Map:
		e.Value = map[string]interface{}(val)
	case marshaler,
		msgp.Extension,
		bool,
		float32,
		float64,
		complex64,
		complex128,
		string,
		[]byte,
		int8,
		int16,
		int32,
		int64,
		int,
		uint,
		uint8,
		uint16,
		uint32,
		uint64,
		time.Time,
		map[string]interface{},
		map[string]string,
		[]interface{}:
		e.Value = val
	default:
		err = ErrInvalidDataType
	}

	return
}

func (e *Entry) Val() Map {
	val, ok := e.Value.(map[string]interface{})
	if ok {
		return Map(val)
	}

	return Map{}
}

func (e *Entry) effective(timeout, now int64) bool {
	if timeout == 0 {
		return true
	}

	if e.UpdatedAt == 0 {
		return now-e.CreatedAt < timeout
	}

	return now-e.UpdatedAt < timeout
}

func (e *Entry) toItem(key string) Item {
	item := Item{
		Key:       key,
		UpdatedAt: e.UpdatedAt,
		UpdateCnt: e.updateCnt,
		HitCnt:    e.hitCnt.Load(),
		MissCnt:   e.missCnt.Load(),
		CreatedAt: e.CreatedAt,
		Value:     fmt.Sprintf("%v", e.Value),
	}

	return item
}

//------------------------ Entry End--------------------------------
