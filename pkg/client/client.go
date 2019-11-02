/*
Copyright © 2019 Dustin Ratcliffe <dustin.k.ratcliffe@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

//Client represents the http client configurables to make http request.
type Client struct {
	BaseURL        *url.URL
	HttpClient     *http.Client
	DefaultHeaders http.Header
	MaxRetry       int
	RetryDelay     time.Duration
}

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

// Config represents the configuration of the api package and all of its structs, methods and functions. These values are
// being read in currently from the app.json secret file.
//
// NOTE: in the future, we are wanting to look into utilizing config maps and secrets to identify between non-sensitive
// and sensitive information.
type Config struct {
	Timeout            time.Duration
	IdleConnTimeout    time.Duration
	Headers            http.Header
	InsecureSkipVerify bool
	MaxConnsPerHost    int
	MaxRetry           int
	RetryDelay         time.Duration
}

// NewClient creates a new instance of a http client.
// These configs values are coming from the Config struct being passed in as a parameter. Once all three are configured,
// we add each struct to our client implemented struct and return it.
func NewClient(conf Config) (*Client, error) {
	// Transport specifies the mechanism by which individual HTTP requests are made
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: conf.InsecureSkipVerify,
		},
		MaxConnsPerHost: conf.MaxConnsPerHost,
		IdleConnTimeout: conf.IdleConnTimeout,
	}

	// Http client
	client := &http.Client{
		Transport: transport,
		Timeout:   conf.Timeout,
	}

	return &Client{HttpClient: client, BaseURL: &url.URL{}, DefaultHeaders: conf.Headers, MaxRetry: conf.MaxRetry, RetryDelay: conf.RetryDelay}, nil
}

// Do determines which http request to make based on the method string value being passed in.
func (c *Client) Do(method string, rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	var request *http.Request
	var err error
	_url := c.BaseURL.ResolveReference(rel)
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

	// Adding headers we configured for our client's headers to our request
	for k, vs := range c.DefaultHeaders {
		for _, v := range vs {
			headers.Add(k, v)
		}
	}
	request.Header = headers
	return c.do(request)
}

// do wraps Client.handle and retries (Client.MaxRetry) with a delay (Client.RetryDelay) if the response code is 5xx or Client.handle returns an error
func (c *Client) do(request *http.Request) (*Response, error) {
	response, err := c.handle(request)
	if (err != nil || IsServerError(response.StatusCode)) && c.MaxRetry > 0 {
		retries := 0
		for retries < c.MaxRetry {
			time.Sleep(c.RetryDelay)

			retries++
			response, err = c.handle(request)
			if err != nil {
				continue
			}

			if !IsServerError(response.StatusCode) {
				return response, err
			}
		}
	}
	return response, err
}

// handle sends an HTTP request and returns an HTTP response, following
// policy (such as redirects, cookies, auth) as configured on the
// client.
//
// An error is returned if caused by client policy (such as
// CheckRedirect), or failure to speak HTTP (such as a network
// connectivity problem). A non-2xx status code doesn't cause an
// error.
//
// If the returned error is nil, the Response will contain a non-empty
// Body. If the Body is not both read to EOF and closed, the Client's
// underlying RoundTripper (typically Transport) may not be able to
// re-use a persistent TCP connection to the server for a subsequent
// "keep-alive" request.
//
// The request Body, if non-nil, will be closed by the underlying
// Transport, even on errors.
//
// On error, any Response can be ignored.
//
// If the server replies with a redirect, the Client first uses the
// CheckRedirect function to determine whether the redirect should be
// followed. If permitted, a 301, 302, or 303 redirect causes
// subsequent requests to use HTTP method GET
// (or HEAD if the original request was HEAD), with no body.
// A 307 or 308 redirect preserves the original HTTP method and body,
// provided that the Request.GetBody function is defined.
// The NewRequest function automatically sets GetBody for common
// standard library body types.
//
// Any returned error will be of type *url.Error. The url.Error
// value's Timeout method will report true if request timed out or was
// canceled.
// handle wraps http.Client.Do
func (c *Client) handle(request *http.Request) (*Response, error) {
	start := time.Now()
	resp, err := c.HttpClient.Do(request)
	if err != nil {
		if err, ok := err.(*url.Error); ok {
			return nil, err.Unwrap()
		} else {
			return nil, err
		}
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

//IsSuccessful checks if code being passed is a successfully code
func IsSuccessful(code int) bool {
	return inRange(code, 200, 300)
}

//IsServerError checks if code being passed is a server error code
func IsServerError(code int) bool {
	return inRange(code, 500, 600)
}

//IsClientError checks if code being passed is a client error code
func IsClientError(code int) bool {
	return inRange(code, 400, 500)
}

//inRange checks if code being passed is between a set of numbers
func InRange(code, a, b int) bool {
	return inRange(code, a, b)
}

//inRange checks if code being passed is between a set of numbers
func inRange(code, a, b int) bool {
	return a <= code && code < b
}
