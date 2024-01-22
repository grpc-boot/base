package elasticsearch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-boot/base/v2/elasticsearch/query"
	"github.com/grpc-boot/base/v2/elasticsearch/result"
	"github.com/grpc-boot/base/v2/utils"
)

func (p *Pool) QueryWithString(ctx context.Context, indexName string, qs *QueryString) (*result.Search, error) {
	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_search", indexName), qs.Marshal(), nil)
	return result.ToSearch(resp, err)
}

type QueryString struct {
	Query       query.StringQuery `json:"query"`
	From        int64             `json:"from"`
	Size        int64             `json:"size,omitempty"`
	Source      any               `json:"_source,omitempty"`
	Sort        query.Sort        `json:"sort,omitempty"`
	SearchAfter []any             `json:"search_after,omitempty"`
}

func NewQueryString(q string) *QueryString {
	qs := &QueryString{
		Size: DefaultPageSize,
	}

	qs.Query = query.StringQuery{
		QueryString: query.String{
			Query: q,
		},
	}

	return qs
}

func (qs *QueryString) SetQuery(q string) *QueryString {
	qs.Query.QueryString.Query = q
	return qs
}

func (qs *QueryString) Marshal() []byte {
	data, _ := utils.JsonMarshal(qs)
	return data
}
