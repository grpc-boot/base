package mq

import "github.com/redis/go-redis/v9"

type Event struct {
	Name string `json:"n"`
	Id   string `json:"i"`
	Data string `json:"d"`
}

func Msg2Event(msg redis.XMessage) Event {
	event := Event{}
	event.Name, _ = msg.Values["n"].(string)
	event.Id, _ = msg.Values["i"].(string)
	event.Data, _ = msg.Values["d"].(string)
	return event
}

func (e Event) AsMsg() redis.XMessage {
	msg := e.ToMsg()
	msg.ID = e.Id
	return msg
}

func (e Event) ToMsg() redis.XMessage {
	return redis.XMessage{
		Values: map[string]interface{}{
			"n": e.Name,
			"i": e.Id,
			"d": e.Data,
		},
	}
}
