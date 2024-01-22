package elasticsearch

import "github.com/grpc-boot/base/v2/elasticsearch/query"

type Query struct {
	From        int64             `json:"from"`
	Size        int64             `json:"size"`
	Source      any               `json:"_source,omitempty"`
	Sort        []query.OrderItem `json:"sort,omitempty"`
	SearchAfter []any             `json:"search_after,omitempty"`
}
