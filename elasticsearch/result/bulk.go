package result

import (
	"errors"

	"github.com/grpc-boot/base/v3/http_client"
)

type Bulk struct {
	Result
	Errors bool                  `json:"errors"`
	Items  []map[string]DocIndex `json:"items"`
}

func (b *Bulk) HasError() bool {
	return b.Result.HasError() || b.Errors
}

func ToBulk(resp *http_client.Response, err error) (*Bulk, error) {
	if err != nil {
		return nil, err
	}

	res := &Bulk{}
	res.Status = resp.GetStatus()
	err = resp.UnmarshalJson(res)
	if err == nil && res.Result.HasError() {
		err = errors.New(res.Error.Reason)
	}
	return res, nil
}
