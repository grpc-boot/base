package base

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var statusPool = sync.Pool{
	New: func() interface{} {
		return &Status{Data: SetValue}
	},
}

// StatusWithCode 指定Code获取一个Status
func StatusWithCode(code codes.Code) *Status {
	sts := statusPool.Get().(*Status)
	sts.Code = code
	if msg, exists := defaultCodeMsg[code]; exists {
		sts.Msg = msg
	}

	return sts
}

// StatusWithCodeMsg 实例化Status
func StatusWithCodeMsg(code codes.Code, msg string) *Status {
	return StatusWithCode(code).WithMsg(msg)
}

// StatusOk 创建携带数据的Status
func StatusOk(data interface{}) *Status {
	return StatusWithCode(OK).WithMsg("ok").WithData(data)
}

// StatusWithJsonUnmarshal 指定json []byte获取一个Status
func StatusWithJsonUnmarshal(data []byte) (*Status, error) {
	sts := &Status{}
	sts.reset()
	if err := JsonUnmarshal(data, sts); err != nil {
		return nil, err
	}
	return sts, nil
}

type Status struct {
	Code codes.Code  `json:"code"`
	Msg  string      `json:"msg"`
	Flag uint8       `json:"flag"`
	Data interface{} `json:"data"`
}

func (s *Status) reset() {
	s.Code = OK
	s.Msg = ""
	s.Data = SetValue
	s.Flag = 0
}

// Close 释放Status到sync.Pool
func (s *Status) Close() {
	s.reset()
	statusPool.Put(s)
}

// IsOK 操作是否OK
func (s *Status) IsOK() bool {
	return s.Is(OK)
}

// Is 判断是否是指定Code
func (s *Status) Is(code codes.Code) bool {
	return s.Code == code
}

// WithMsg 附加自定义message
func (s *Status) WithMsg(msg string) *Status {
	s.Msg = msg
	return s
}

// WithFlag 附加flag
func (s *Status) WithFlag(flag uint8) *Status {
	s.Flag = flag
	return s
}

// WithData 附加自定义Data
func (s *Status) WithData(data interface{}) *Status {
	s.Data = data
	return s
}

// Error 转换为error
func (s *Status) Error() error {
	return s.ToError()
}

// JsonMarshal _
func (s *Status) JsonMarshal() []byte {
	data, _ := JsonMarshal(s)
	return data
}

// ConvertGrpcStatus 转换为grpc状态码
func (s *Status) ConvertGrpcStatus(details ...proto.Message) (*status.Status, error) {
	return status.New(s.Code, s.Msg).WithDetails(details...)
}

func (s *Status) ToError() error {
	if s.Code == OK {
		return nil
	}

	return NewError(s.Code, s.Msg)
}
