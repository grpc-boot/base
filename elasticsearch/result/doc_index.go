package result

import (
	"github.com/grpc-boot/base/v2/http_client"
)

type DocIndex struct {
	Result

	Id          string `json:"_id"`
	Version     int64  `json:"_version"`
	Res         string `json:"result"`
	SeqNo       int64  `json:"_seq_no"`
	PrimaryTerm int64  `json:"_primary_term"`
}

func ToDocIndex(resp *http_client.Response, err error) (*DocIndex, error) {
	if err != nil {
		return nil, err
	}

	res := &DocIndex{}
	res.Status = resp.GetStatus()
	err = resp.UnmarshalJson(res)
	return res, err
}
