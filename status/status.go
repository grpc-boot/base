package status

import (
	"errors"
	"fmt"

	"github.com/grpc-boot/base/v3/utils"
)

type Status struct {
	Code    uint32 `json:"code" yaml:"code"`
	Message string `json:"msg" yaml:"msg"`
	Data    any    `json:"data" yaml:"data"`
	Flag    uint16 `json:"flag" yaml:"flag"`
}

func JsonUnmarshalStatus(data []byte) (*Status, error) {
	var status Status
	if err := utils.JsonUnmarshal(data, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func StatusOk(data any) *Status {
	return &Status{
		Code:    Ok,
		Message: "ok",
		Data:    data,
	}
}

func StatusError(code uint32, flag uint16) *Status {
	return &Status{
		Code: code,
		Flag: flag,
	}
}

func (s *Status) IsCode(code uint32) bool {
	return s.Code == code
}

func (s *Status) WithMsg(msg string) *Status {
	s.Message = msg
	return s
}

func (s *Status) WithFlag(flag uint16) *Status {
	s.Flag = flag
	return s
}

func (s *Status) Error() string {
	return fmt.Sprintf("%s: [code: %d, flag: %d]", s.Message, s.Code, s.Flag)
}

func (s *Status) ToError() error {
	return errors.New(s.Error())
}

func (s *Status) Reset() {
	s.Code = Ok
	s.Message = ""
	s.Data = nil
	s.Flag = 0
}
