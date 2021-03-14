package rpcx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/refx"
)

type HttpClientOptions struct {
	DialTimeout         time.Duration `dft:"3s"`
	Timeout             time.Duration `dft:"6s"`
	MaxIdleConnsPerHost int           `dft:"2"`
}

func NewHttpClient() *HttpClient {
	var options HttpClientOptions
	refx.SetDefaultValueP(&options)

	return NewHttpClientWithOptions(&options)
}

func NewHttpClientWithOptions(options *HttpClientOptions) *HttpClient {
	return &HttpClient{
		options: options,
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, options.DialTimeout)
					if err != nil {
						return nil, err
					}
					return c, nil
				},
				MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
			},
			Timeout: options.Timeout,
		},
	}
}

type HttpClient struct {
	options    *HttpClientOptions
	httpClient *http.Client
}

type HttpError ErrorDetail

func (e *HttpError) Error() string {
	return fmt.Sprintf("%v. status: [%v], code: [%v], requestID: [%v], refer: [%v]", e.Message, e.Status, e.Code, e.RequestID, e.Refer)
}

func (c *HttpClient) Get(uri string, params map[string]string, reqMeta map[string]string, req interface{}, resMeta *map[string]string, res interface{}) error {
	return c.Do("GET", uri, params, reqMeta, req, resMeta, res)
}

func (c *HttpClient) Post(uri string, params map[string]string, reqMeta map[string]string, req interface{}, resMeta *map[string]string, res interface{}) error {
	return c.Do("POST", uri, params, reqMeta, req, resMeta, res)
}

func (c *HttpClient) Do(method string, uri string, params map[string]string, reqMeta map[string]string, req interface{}, resMeta *map[string]string, res interface{}) error {
	var r io.Reader
	if req != nil {
		buf, err := json.Marshal(req)
		if err != nil {
			return errors.Wrap(err, "json.Marshal err failed")
		}
		r = bytes.NewReader(buf)
	}

	hreq, err := http.NewRequest(method, uri, r)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
	}
	for key, val := range reqMeta {
		hreq.Header.Set(key, val)
	}
	if params != nil {
		q := hreq.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		hreq.URL.RawQuery = q.Encode()
	}

	hres, err := c.httpClient.Do(hreq)
	if err != nil {
		return errors.Wrap(err, "client.Do failed")
	}
	defer hres.Body.Close()

	if resMeta != nil {
		for key, val := range hres.Header {
			(*resMeta)[key] = val[0]
		}
	}

	buf, err := ioutil.ReadAll(hres.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll failed")
	}

	if hres.StatusCode != http.StatusOK {
		var herr HttpError
		if err := json.Unmarshal(buf, &herr); err != nil {
			return &HttpError{
				Status:  int32(hres.StatusCode),
				Code:    hres.Status,
				Message: string(buf),
			}
		}
		return &herr
	}

	if len(buf) != 0 {
		if err := json.Unmarshal(buf, res); err != nil {
			return errors.Wrap(err, "json.Unmarshal failed")
		}
	}

	return nil
}
