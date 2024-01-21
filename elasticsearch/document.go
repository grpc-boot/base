package elasticsearch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-boot/base/v2/elasticsearch/result"
)

func (p *Pool) DocIndex(ctx context.Context, indexName string, body Body) (res *result.DocIndex, err error) {
	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_doc", indexName), body.Marshal(), nil)
	return result.ToDocIndex(resp, err)
}

func (p *Pool) DocIndexWithId(ctx context.Context, indexName, id string, body Body) (res *result.DocIndex, err error) {
	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_doc/%s", indexName, id), body.Marshal(), nil)
	return result.ToDocIndex(resp, err)
}

func (p *Pool) DocMGet(ctx context.Context, indexName string, idList ...string) (res *result.Documents, err error) {
	body := make(Body, 1)
	body.WithArgs(WithArg("ids", idList))

	resp, err := p.Request(ctx, http.MethodGet, fmt.Sprintf("%s/_mget", indexName), body.Marshal(), nil)

	return result.ToDocuments(resp, err)
}

func (p *Pool) DocDel(ctx context.Context, indexName, id string) (res *result.DocIndex, err error) {
	resp, err := p.Request(ctx, http.MethodDelete, fmt.Sprintf("%s/_doc/%s", indexName, id), nil, nil)

	return result.ToDocIndex(resp, err)
}

func (p *Pool) DocUpdate(ctx context.Context, indexName, id string, setter Setter) (res *result.DocIndex, err error) {
	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_update/%s", indexName, id), setter.ToBody(), nil)
	return result.ToDocIndex(resp, err)
}

func (p *Pool) DocUpdateWithOptimistic(ctx context.Context, indexName, id string, setter Setter, seqNo, primaryTerm int64) (res *result.DocIndex, err error) {
	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_update/%s?if_seq_no=%d&if_primary_term=%d", indexName, id, seqNo, primaryTerm), setter.ToBody(), nil)
	return result.ToDocIndex(resp, err)
}

func (p *Pool) DocFieldIncr(ctx context.Context, indexName, id, field string, increment int64) (res *result.DocIndex, err error) {
	body := fmt.Sprintf(`{
		"script": {
			"source": "ctx._source.%s += params.count",
			"lang": "painless",
			"params" : {"count" : %d}
		}
	}`, field, increment)

	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("%s/_update/%s", indexName, id), []byte(body), nil)
	return result.ToDocIndex(resp, err)
}
