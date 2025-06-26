package http_client

import (
	"github.com/grpc-boot/base/v3/utils"
)

type Response struct {
	status int
	body   []byte
}

func NewResp(status int) *Response {
	return &Response{
		status: status,
	}
}

func (rp *Response) Is(status int) bool {
	return rp.status == status
}

func (rp *Response) GetStatus() int {
	return rp.status
}

func (rp *Response) SetBody(body []byte) *Response {
	rp.body = body
	return rp
}

func (rp *Response) GetBody() []byte {
	return rp.body
}

func (rp *Response) UnmarshalJson(out any) error {
	return utils.JsonUnmarshal(rp.body, out)
}
