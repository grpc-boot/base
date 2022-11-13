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
	Key         string `json:"key" msg:"key"`
	Value       []byte `json:"value" msg:"value"`
	Hit         int64  `json:"hit" msg:"hit"`
	Miss        int64  `json:"miss" msg:"miss"`
	InvokeCount int64  `json:"invokeCount" msg:"invokeCount"`
	CreatedAt   int64  `json:"createdAt" msg:"createdAt"`
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
