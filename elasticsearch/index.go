package elasticsearch

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-boot/base/v2/elasticsearch/result"
	"github.com/grpc-boot/base/v2/http_client"
)

func (p *Pool) Index(ctx context.Context, name string, args ...Arg) (res *result.Index, err error) {
	var resp *http_client.Response

	if len(args) == 0 {
		resp, err = p.Request(ctx, http.MethodPut, fmt.Sprintf("%s", name), nil, nil)
	} else {
		body := make(Body, len(args))
		body.WithArgs(args...)
		resp, err = p.Request(ctx, http.MethodPut, name, body.Marshal(), nil)
	}

	return result.ToIndex(resp, err)
}

func (p *Pool) IndexDel(ctx context.Context, name string) (res *result.Index, err error) {
	resp, err := p.Request(ctx, http.MethodDelete, name, nil, nil)
	if err != nil {
		return
	}

	return result.ToIndex(resp, err)
}

func (p *Pool) IndexSetting(ctx context.Context, name string, args ...Arg) (res *result.Index, err error) {
	if len(args) == 0 {
		err = ErrArgsEmpty
		return
	}

	body := make(Body, len(args))
	body.WithArgs(args...)

	resp, err := p.Request(ctx, http.MethodPut, fmt.Sprintf("%s/_settings", name), body.Marshal(), nil)

	return result.ToIndex(resp, err)
}

func (p *Pool) IndexSettingGet(ctx context.Context, target string) (response *http_client.Response, err error) {
	return p.Request(ctx, http.MethodGet, fmt.Sprintf("%s/_settings", target), nil, nil)
}

func (p *Pool) IndexMapping(ctx context.Context, name string, properties result.MappingProperties) (res *result.Index, err error) {
	if len(properties) == 0 {
		err = ErrPropertiesEmpty
		return
	}

	body := make(Body, 1)
	body.WithProperties(properties)

	resp, err := p.Request(ctx, http.MethodPut, fmt.Sprintf("%s/_mapping", name), body.Marshal(), nil)

	return result.ToIndex(resp, err)
}

func (p *Pool) IndexMappingGet(ctx context.Context, target string) (res *result.IndexMapping, err error) {
	var resp *http_client.Response
	if target == "" {
		resp, err = p.Request(ctx, http.MethodGet, "_mapping", nil, nil)
	} else {
		resp, err = p.Request(ctx, http.MethodGet, fmt.Sprintf("%s/_mapping", target), nil, nil)
	}

	return result.ToMapping(resp, err)
}
