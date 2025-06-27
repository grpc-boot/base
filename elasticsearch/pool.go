package elasticsearch

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-boot/base/v3/elasticsearch/result"
	"github.com/grpc-boot/base/v3/http_client"
	"github.com/grpc-boot/base/v3/utils"
)

type Pool struct {
	pool      *http_client.Pool
	opt       Options
	basicAuth string
}

func NewPool(opt Options) (pool *Pool) {
	opt.BaseUrl = strings.TrimSuffix(opt.BaseUrl, "/")

	pool = &Pool{
		pool: http_client.NewPool(opt.httpOptions()),
		opt:  opt,
	}

	if opt.UserName != "" {
		auth := opt.UserName + ":" + opt.Password
		pool.basicAuth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(utils.String2Bytes(auth)))
	}
	return
}

func (p *Pool) Options() Options {
	return p.opt
}

func (p *Pool) Request(ctx context.Context, method, path string, body []byte, headers http_client.Headers) (response *http_client.Response, err error) {
	var (
		url = fmt.Sprintf("%s/%s", p.opt.BaseUrl, path)
	)

	if len(headers) == 0 {
		headers = http_client.Headers{
			"Content-Type": "application/vnd.elasticsearch+json; compatible-with=7",
			"Accept":       "application/vnd.elasticsearch+json; compatible-with=7",
		}
	}

	if p.basicAuth != "" {
		headers["Authorization"] = p.basicAuth
	}

	return p.pool.Request(ctx, method, url, body, headers)
}

func (p *Pool) SearchBySql(ctx context.Context, size int64, format, sqlStr string, params []any, bodyArgs ...Arg) (res *result.Sql, err error) {
	if format == "" {
		format = "json"
	}

	if size == 0 {
		size = DefaultPageSize
	}

	body := make(Body, 3+len(bodyArgs))
	body.WithSql(sqlStr, params...)
	body.WithFetchSize(size)
	body.WithArgs(bodyArgs...)

	reqBody, err := body.Marshal()
	if err != nil {
		return nil, err
	}
	body = nil

	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("_sql?format=%s", format), reqBody, nil)
	return result.ToSql(resp, err)
}

func (p *Pool) QueryWithCursor(ctx context.Context, cursor, format string, args ...Arg) (res *result.Sql, err error) {
	if format == "" {
		format = "json"
	}

	body := make(Body, 1+len(args))
	body.WithCursor(cursor)
	body.WithArgs(args...)

	reqBody, err := body.Marshal()
	if err != nil {
		return nil, err
	}
	body = nil

	resp, err := p.Request(ctx, http.MethodPost, fmt.Sprintf("_sql?format=%s", format), reqBody, nil)
	return result.ToSql(resp, err)
}

func (p *Pool) CloseSqlCursor(ctx context.Context, cursor string) (response *http_client.Response, err error) {
	body := make(Body, 1)
	body.WithCursor(cursor)

	reqBody, err := body.Marshal()
	if err != nil {
		return nil, err
	}
	body = nil

	return p.Request(ctx, http.MethodPost, "_sql/close", reqBody, nil)
}

func (p *Pool) Bulk(ctx context.Context, bi *BulkItem) (res *result.Bulk, err error) {
	resp, err := p.Request(ctx, http.MethodPost, "_bulk", bi.Marshal(), nil)
	return result.ToBulk(resp, err)
}
