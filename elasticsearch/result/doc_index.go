package result

import (
	"errors"

	"github.com/grpc-boot/base/v3/http_client"
)

type DocIndex struct {
	Result

	Id          string `json:"_id"`
	Version     int64  `json:"_version"`
	Res         string `json:"result"`
	SeqNo       int64  `json:"_seq_no"`
	PrimaryTerm int64  `json:"_primary_term"`
}

func (di *DocIndex) IsRes(result string) bool {
	return di.Res == result
}

func ToDocIndex(resp *http_client.Response, err error) (*DocIndex, error) {
	if err != nil {
		return nil, err
	}

	res := &DocIndex{}
	res.Status = resp.GetStatus()
	err = resp.UnmarshalJson(res)
	if err == nil && res.HasError() {
		err = errors.New(res.Error.Reason)
	}

	return res, err
}
