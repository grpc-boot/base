package base

import jsoniter "github.com/json-iterator/go"

type Msg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (m *Msg) Json() (msg []byte) {
	msg, _ = jsoniter.Marshal(m)
	return
}
