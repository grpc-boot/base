package result

import (
	"errors"

	"github.com/grpc-boot/base/v2/http_client"
)

type Index struct {
	Result

	Acknowledged bool `json:"acknowledged"`
}

func (i *Index) Ok() bool {
	return i.Acknowledged
}

func ToIndex(resp *http_client.Response, err error) (*Index, error) {
	if err != nil {
		return nil, err
	}

	res := &Index{}
	res.Status = resp.GetStatus()

	err = resp.UnmarshalJson(res)
	if err == nil && res.HasError() {
		err = errors.New(res.Error.Reason)
	}

	return res, err
}
