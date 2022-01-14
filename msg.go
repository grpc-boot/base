package base

// Msg 通用消息
type Msg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Json 转换为Json
func (m *Msg) Json() (msg []byte) {
	msg, _ = JsonEncode(m)
	return
}
