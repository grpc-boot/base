package shardmap

type Event uint8

const (
	Create Event = 1
	Update Event = 2
	Delete Event = 3
)

// ChangeEvent 修改事件
type ChangeEvent struct {
	Type     Event
	Key      interface{}
	OldValue interface{}
	Value    interface{}
}
