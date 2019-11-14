package client

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client ...
type Client struct {
	client     *http.Client
	url        *url.URL
	header     http.Header
	maxRetry   int
	retryDelay time.Duration
}

// Response ...
type Response struct {
	Status           string // e.g. "200 OK"
	StatusCode       int    // e.g. 200
	Proto            string // e.g. "HTTP/1.0"
	ProtoMajor       int    // e.g. 1
	ProtoMinor       int    // e.g. 0
	Header           http.Header
	Body             []byte
	ContentLength    int64
	TransferEncoding []string
	Uncompressed     bool
	Trailer          http.Header
	Duration         time.Duration
}

// Config ...
type Config struct {
	Timeout            time.Duration
	IdleConnTimeout    time.Duration
	Headers            http.Header
	InsecureSkipVerify bool
	MaxConnsPerHost    int
	maxRetry           int
	retryDelay         time.Duration
}

// NewClient ...
func NewClient(conf Config) (*Client, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: conf.InsecureSkipVerify,
		},
		MaxConnsPerHost: conf.MaxConnsPerHost,
		IdleConnTimeout: conf.IdleConnTimeout,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   conf.Timeout,
	}

	return &Client{client: client, url: &url.URL{}, header: conf.Headers, maxRetry: conf.maxRetry, retryDelay: conf.retryDelay}, nil
}

// Do ...
func (c *Client) Do(method string, rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	var request *http.Request
	var err error
	_url := c.url.ResolveReference(rel)
	switch strings.ToUpper(method) {
	case http.MethodGet:
		request, err = http.NewRequest(http.MethodGet, _url.String(), nil)
		if err != nil {
			return nil, err
		}
	case http.MethodHead:
		request, err = http.NewRequest(http.MethodHead, _url.String(), nil)
		if err != nil {
			return nil, err
		}
	case http.MethodPost:
		request, err = http.NewRequest(http.MethodPost, _url.String(), body)
		if err != nil {
			return nil, err
		}
	case http.MethodPut:
		request, err = http.NewRequest(http.MethodPut, _url.String(), body)
		if err != nil {
			return nil, err
		}
	case http.MethodPatch:
		request, err = http.NewRequest(http.MethodPatch, _url.String(), body)
		if err != nil {
			return nil, err
		}
	case http.MethodDelete:
		request, err = http.NewRequest(http.MethodDelete, _url.String(), nil)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported method")
	}

	for k, vs := range c.header {
		for _, v := range vs {
			headers.Add(k, v)
		}
	}
	request.Header = headers
	return c.do(request)
}

func (c *Client) do(request *http.Request) (*Response, error) {
	response, err := c.handle(request)
	if (err != nil || inRange(response.StatusCode, 500, 600)) && c.maxRetry > 0 {
		retries := 0
		for retries < c.maxRetry {
			time.Sleep(c.retryDelay)

			retries++
			response, err = c.handle(request)
			if err != nil {
				continue
			}

			if !inRange(response.StatusCode, 500, 600) {
				return response, err
			}
		}
	}
	return response, err
}

func (c *Client) handle(request *http.Request) (*Response, error) {
	start := time.Now()
	resp, err := c.client.Do(request)
	if err != nil {
		if err, ok := err.(*url.Error); ok {
			return nil, err.Unwrap()
		}
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := Response{
		Status:           resp.Status,
		StatusCode:       resp.StatusCode,
		Proto:            resp.Proto,
		ProtoMajor:       resp.ProtoMajor,
		ProtoMinor:       resp.ProtoMinor,
		Header:           resp.Header,
		Body:             body,
		ContentLength:    resp.ContentLength,
		TransferEncoding: resp.TransferEncoding,
		Uncompressed:     resp.Uncompressed,
		Trailer:          resp.Trailer,
		Duration:         time.Now().Sub(start),
	}
	return &response, nil
}

// InRange ...
func InRange(code, a, b int) bool {
	return inRange(code, a, b)
}

func inRange(code, a, b int) bool {
	return a <= code && code < b
}
