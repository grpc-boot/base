package http_client

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/grpc-boot/base/v2/logger"
	"github.com/grpc-boot/base/v2/utils"

	"go.uber.org/zap/zapcore"
)

type Client struct {
	client *http.Client
	opt    *Options
}

func NewClient(opt Options) *Client {
	c := &Client{
		opt: &opt,
	}

	c.client = &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   c.opt.DisableKeepAlives,
			DisableCompression:  c.opt.DisableCompression,
			MaxIdleConns:        c.opt.MaxIdleConns,
			MaxIdleConnsPerHost: c.opt.MaxIdleConnsPerHost,
			MaxConnsPerHost:     c.opt.MaxConnsPerHost,
			IdleConnTimeout:     c.opt.IdleConnTimeout(),
			WriteBufferSize:     c.opt.WriteBufferSize,
			ReadBufferSize:      c.opt.ReadBufferSize,
		},
		Timeout: c.opt.Timeout(),
	}

	return c
}

func (c *Client) OptionsTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodOptions, url, nil, headers)
}

func (c *Client) Options(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodOptions, url, nil, headers)
}

func (c *Client) HeadTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodHead, url, nil, headers)
}

func (c *Client) Head(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodHead, url, nil, headers)
}

func (c *Client) GetTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodGet, url, nil, headers)
}

func (c *Client) Get(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodGet, url, nil, headers)
}

func (c *Client) PostTimeout(timeout time.Duration, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodPost, url, body, headers)
}

func (c *Client) Post(ctx context.Context, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodPost, url, body, headers)
}

func (c *Client) PutTimeout(timeout time.Duration, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodPut, url, body, headers)
}

func (c *Client) Put(ctx context.Context, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodPut, url, body, headers)
}

func (c *Client) DeleteTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodDelete, url, nil, headers)
}

func (c *Client) Delete(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodDelete, url, nil, headers)
}

func (c *Client) RequestTimeout(timeout time.Duration, method, url string, body []byte, headers Headers) (rp *Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.Request(ctx, method, url, body, headers)
}

func (c *Client) Request(ctx context.Context, method, url string, body []byte, headers Headers) (rp *Response, err error) {
	var (
		start = time.Now()
		req   *http.Request
	)

	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("http request begin",
			logger.Uri(url),
			logger.Method(method),
			logger.Headers(headers),
			logger.Params(utils.Bytes2String(body)),
		)
	}

	if len(body) > 0 {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	if err != nil {
		return
	}

	if len(headers) > 0 {
		for key, value := range headers {
			if key == "Host" {
				req.Host = value
				continue
			}

			req.Header.Set(key, value)
		}
	}

	rp, err = c.Do(req)

	if logger.IsLevel(zapcore.DebugLevel) {
		logger.ZapDebug("http request done",
			logger.Uri(url),
			logger.Method(method),
			logger.Headers(headers),
			logger.Params(utils.Bytes2String(body)),
			logger.Duration(time.Since(start)),
		)
	}

	if err != nil {
		logger.ZapDebug("http request",
			logger.Error(err),
			logger.Uri(url),
			logger.Method(method),
			logger.Headers(headers),
			logger.Params(utils.Bytes2String(body)),
			logger.Duration(time.Since(start)),
		)
	}

	return
}

func (c *Client) Do(req *http.Request) (rp *Response, err error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	rp = NewResp(resp.StatusCode)
	rp.body, err = io.ReadAll(resp.Body)

	return
}
