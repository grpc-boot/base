//go:generate msgp
package cache

import (
	"time"
)

type Bucket struct {
	Items map[string]*Item `json:"items" msg:"items"`
}

type Item struct {
	_lock       int64
	Key         string `json:"key" msg:"-"`
	Value       []byte `json:"value" msg:"value"`
	Hit         int64  `json:"hit" msg:"-"`
	Miss        int64  `json:"miss" msg:"-"`
	InvokeCount int64  `json:"invokeCount" msg:"-"`
	CreatedAt   int64  `json:"createdAt" msg:"createdAt"`
}

func (i *Item) save(value []byte) {
	i.CreatedAt = time.Now().Unix()
	i.Value = value
	i.InvokeCount++
}

// effective 数据是否有效
func (i *Item) effective(timeoutSecond int64) (ok bool) {
	ok = i.CreatedAt+timeoutSecond > time.Now().Unix()
	if ok {
		i.Hit++
	} else {
		i.Miss++
	}

	return ok
}
