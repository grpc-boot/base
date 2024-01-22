package result

import (
	"errors"

	"github.com/grpc-boot/base/v2/http_client"
)

const (
	RelationEq = `eq`
	RelationGt = `gt`
)

type Search struct {
	Result

	Hits     Hits  `json:"hits"`
	TimedOut bool  `json:"timed_out"`
	Took     int64 `json:"took"`
}

type Hits struct {
	Total    Total     `json:"total"`
	MaxScore any       `json:"max_score"`
	Hits     []HitItem `json:"hits"`
}

type Total struct {
	Value    int64  `json:"value"`
	Relation string `json:"relation"`
}

type HitItem struct {
	Index  string `json:"_index"`
	Id     string `json:"_id"`
	Score  any    `json:"_score"`
	Source Source `json:"_source"`
	Sort   []any  `json:"sort"`
}

func ToSearch(resp *http_client.Response, err error) (*Search, error) {
	if err != nil {
		return nil, err
	}

	res := &Search{}
	res.Status = resp.GetStatus()

	err = resp.UnmarshalJson(res)
	if err == nil && res.HasError() {
		err = errors.New(res.Error.Reason)
	}

	return res, err
}
