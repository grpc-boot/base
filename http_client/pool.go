package http_client

import (
	"bytes"
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/grpc-boot/base/v3/logger"
	"github.com/grpc-boot/base/v3/utils"

	"go.uber.org/zap"
)

type Pool struct {
	client *http.Client
	opt    *Options
}

func NewPool(opt Options) *Pool {
	transport := &http.Transport{
		Proxy:             http.ProxyFromEnvironment,
		ForceAttemptHTTP2: true,
		DialContext: (&net.Dialer{
			Timeout:   opt.DialTimeout(),
			KeepAlive: opt.KeepaliveTime(),
		}).DialContext,
		MaxIdleConns:        opt.MaxIdleConns,
		MaxIdleConnsPerHost: opt.MaxIdleConnsPerHost,
		MaxConnsPerHost:     opt.MaxConnsPerHost,
		IdleConnTimeout:     opt.IdleConnTimeout(),
		WriteBufferSize:     opt.WriteBufferSize,
		ReadBufferSize:      opt.ReadBufferSize,
		DisableKeepAlives:   opt.DisableKeepAlives,
		DisableCompression:  opt.DisableCompression,
	}

	return &Pool{
		opt: &opt,
		client: &http.Client{
			Transport: transport,
			Timeout:   opt.Timeout(),
		},
	}
}

func (c *Pool) OptionsTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodOptions, url, nil, headers)
}

func (c *Pool) Options(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodOptions, url, nil, headers)
}

func (c *Pool) HeadTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodHead, url, nil, headers)
}

func (c *Pool) Head(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodHead, url, nil, headers)
}

func (c *Pool) GetTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodGet, url, nil, headers)
}

func (c *Pool) Get(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodGet, url, nil, headers)
}

func (c *Pool) PostTimeout(timeout time.Duration, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodPost, url, body, headers)
}

func (c *Pool) Post(ctx context.Context, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodPost, url, body, headers)
}

func (c *Pool) PutTimeout(timeout time.Duration, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodPut, url, body, headers)
}

func (c *Pool) Put(ctx context.Context, url string, body []byte, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodPut, url, body, headers)
}

func (c *Pool) DeleteTimeout(timeout time.Duration, url string, headers Headers) (rp *Response, err error) {
	return c.RequestTimeout(timeout, http.MethodDelete, url, nil, headers)
}

func (c *Pool) Delete(ctx context.Context, url string, headers Headers) (rp *Response, err error) {
	return c.Request(ctx, http.MethodDelete, url, nil, headers)
}

func (c *Pool) RequestTimeout(timeout time.Duration, method, url string, body []byte, headers Headers) (rp *Response, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.Request(ctx, method, url, body, headers)
}

func (c *Pool) Request(ctx context.Context, method, url string, body []byte, headers Headers) (rp *Response, err error) {
	var (
		start = time.Now()
		req   *http.Request
	)

	if logger.IsDebug() {
		logger.Debug("http request begin",
			zap.String("Url", url),
			zap.String("Method", method),
			zap.Any("Headers", headers),
			zap.String("Body", utils.Bytes2String(body)),
		)
	}

	if len(body) > 0 {
		req, err = http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}

	if err != nil {
		logger.Error("create http request failed",
			zap.String("Url", url),
			zap.String("Method", method),
			zap.Any("Headers", headers),
			zap.String("Body", utils.Bytes2String(body)),
			zap.Duration("Duration", time.Since(start)),
			zap.NamedError("Error", err),
		)
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
	if err != nil {
		logger.Error("http request failed",
			zap.String("Url", url),
			zap.String("Method", method),
			zap.Any("Headers", headers),
			zap.String("Body", utils.Bytes2String(body)),
			zap.Duration("Duration", time.Since(start)),
			zap.NamedError("Error", err),
		)
	} else {
		if logger.IsDebug() {
			logger.Debug("http request done",
				zap.String("Url", url),
				zap.String("Method", method),
				zap.Any("Headers", headers),
				zap.String("Body", utils.Bytes2String(body)),
				zap.Duration("Duration", time.Since(start)),
			)
		}
	}

	return
}

func (c *Pool) Do(req *http.Request) (rp *Response, err error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	rp = NewResp(resp.StatusCode)
	rp.body, err = io.ReadAll(resp.Body)
	return
}
