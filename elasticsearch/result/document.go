package result

import (
	"errors"

	"github.com/grpc-boot/base/v3/http_client"
)

type Documents struct {
	Result

	Docs []Document `json:"docs"`
}

type Document struct {
	DocIndex

	Source Source `json:"_source"`
}

func ToDocuments(resp *http_client.Response, err error) (*Documents, error) {
	if err != nil {
		return nil, err
	}

	res := &Documents{}
	res.Status = resp.GetStatus()
	err = resp.UnmarshalJson(res)
	if err == nil && res.HasError() {
		err = errors.New(res.Error.Reason)
	}
	return res, err
}
