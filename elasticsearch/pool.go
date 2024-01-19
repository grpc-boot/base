package elasticsearch

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-boot/base/v2/http_client"
	"github.com/grpc-boot/base/v2/internal"
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
		pool.basicAuth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(internal.String2Bytes(auth)))
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

func (p *Pool) SearchBySql(ctx context.Context, size int64, format, sqlStr string, params []any, bodyArgs ...Arg) (response *http_client.Response, err error) {
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

	return p.Request(ctx, http.MethodPost, fmt.Sprintf("_sql?format=%s", format), body.Marshal(), nil)
}

func (p *Pool) Query(ctx context.Context, query *Query, format string, args ...Arg) (response *http_client.Response, err error) {
	querySql, params := query.Sql()
	return p.SearchBySql(ctx, query.limit, format, querySql, params, args...)
}

func (p *Pool) QueryWithCursor(ctx context.Context, cursor, format string, args ...Arg) (response *http_client.Response, err error) {
	if format == "" {
		format = "json"
	}

	body := make(Body, 1+len(args))
	body.WithCursor(cursor)
	body.WithArgs(args...)

	return p.Request(ctx, http.MethodPost, fmt.Sprintf("_sql?format=%s", format), body.Marshal(), nil)
}

func (p *Pool) CloseSqlCursor(ctx context.Context, cursor string) (response *http_client.Response, err error) {
	body := make(Body, 1)
	body.WithCursor(cursor)

	return p.Request(ctx, http.MethodPost, "_sql/close", body.Marshal(), nil)
}
