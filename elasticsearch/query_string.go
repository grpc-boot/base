package elasticsearch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-boot/base/v3/elasticsearch/result"
	"github.com/grpc-boot/base/v3/utils"
)

func (p *Pool) QueryWithString(ctx context.Context, indexName string, qs *QueryString) (*result.Search, error) {
	body, err := utils.JsonMarshal(qs)
	if err != nil {
		return nil, err
	}

	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_search", indexName), body, nil)
	return result.ToSearch(resp, err)
}

type QueryString struct {
	Query       StringQuery `json:"query"`
	From        int64       `json:"from"`
	Size        int64       `json:"size,omitempty"`
	Source      any         `json:"_source,omitempty"`
	Sort        Sort        `json:"sort,omitempty"`
	SearchAfter []any       `json:"search_after,omitempty"`
}

func NewQueryString(q string) *QueryString {
	qs := &QueryString{
		Size: DefaultPageSize,
	}

	qs.Query = StringQuery{
		QueryString: String{
			Query: q,
		},
	}

	return qs
}

func (qs *QueryString) SetQuery(q string) *QueryString {
	qs.Query.QueryString.Query = q
	return qs
}
