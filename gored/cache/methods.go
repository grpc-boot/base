package cache

import (
	"time"

	"github.com/grpc-boot/base/v2/kind/msg"
)

func (i *Item) Marshal() (data []byte, err error) {
	data = make([]byte, 0, i.Msgsize())
	return i.MarshalMsg(data)
}

func (i *Item) Hit(timeoutSeconds int64, current int64) bool {
	return i.UpdatedAt+timeoutSeconds > current
}

func (i *Item) PackValue() {
	switch val := i.Value.(type) {
	case msg.Map:
		i.Value = map[string]interface{}(val)
	}
}

func (i *Item) Bool() bool {
	return msg.Bool(i.Value)
}

func (i *Item) Int() int64 {
	return msg.Int(i.Value)
}

func (i *Item) Uint() uint64 {
	return msg.Uint(i.Value)
}

func (i *Item) Float() float64 {
	return msg.Float(i.Value)
}

func (i *Item) String() string {
	return msg.String(i.Value)
}

func (i *Item) Bytes() []byte {
	return msg.Bytes(i.Value)
}

func (i *Item) Slice() []interface{} {
	return msg.Slice(i.Value)
}

func (i *Item) Map() msg.Map {
	return msg.MsgMap(i.Value)
}

func (i *Item) StringMap() map[string]string {
	return msg.StringMap(i.Value)
}

func (i *Item) Time() time.Time {
	return msg.Time(i.Value)
}
