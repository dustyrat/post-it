package client

import (
	"crypto/tls"
	"errors"
	log "github.com/sirupsen/logrus"
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
}

//Config represents the configuration of the api package and all of its structs, methods and functions. These values are
// being read in currently from the app.json secret file.
//
//NOTE: in the future, we are wanting to look into utilizing config maps and secrets to identify between non-sensitive
// and sensitive information.
type Config struct {
	URL string `json:"url"`

	Timeout         time.Duration `json:"timeout-ms"`
	IdleConnTimeout time.Duration `json:"idle-connection-timeout-ms"`

	Headers            map[string]string `json:"default-headers"`
	InsecureSkipVerify bool              `json:"insecure-skip-verify"`

	MaxConnsPerHost     int `json:"max-connection-per-host"`
	MaxIdleConns        int `json:"max-idle-connections"`
	MaxIdleConnsPerHost int `json:"max-idle-connections-per-host"`
	MaxRetry            int `json:"max-retry"`

	RetryDelay time.Duration `json:"retry-delay-ms"`
}

//NewClient creates a new instance of a http client.
//In order to create a new client, we need to instantiate and configure three structs from the http package:
//	- transport{}
//	- client{}
//	- headers{}
//
//These configs values are coming from the Config struct being passed in as a parameter. Once all three are configured,
// we add each struct to our client implemented struct and return it.
func NewClient(conf Config) (*Client, error) {
	_url, err := url.Parse(conf.URL)
	if err != nil {
		return nil, err
	}

	// Transport specifies the mechanism by which individual HTTP requests are made
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: conf.InsecureSkipVerify,
		},
		MaxConnsPerHost:     conf.MaxConnsPerHost,
		MaxIdleConns:        conf.MaxIdleConns,
		MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		IdleConnTimeout:     conf.IdleConnTimeout * time.Millisecond,
	}

	// Http client
	client := &http.Client{
		Transport: transport,
		Timeout:   conf.Timeout * time.Millisecond,
	}

	// A Header represents the key-value pairs in an HTTP header. We are creating all string key pairs headers that are
	// being read in from the config file.
	headers := http.Header{}
	for k, v := range conf.Headers {
		headers.Set(k, v)
	}

	return &Client{HttpClient: client, BaseURL: _url, DefaultHeaders: headers, MaxRetry: conf.MaxRetry, RetryDelay: conf.RetryDelay * time.Millisecond}, nil
}

//Do determines which http request to make based on the method string value being passed in. If a valid method value, it
// will pass the the remaining parameters to its respected request:
// 	- rel (relationship url)
// 	- headers (http headers)
// 	- body (request payload)
func (c *Client) Do(method string, rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	switch strings.ToUpper(method) {
	case http.MethodGet:
		return c.Get(rel, headers)
	case http.MethodPut:
		return c.Put(rel, headers, body)
	case http.MethodPost:
		return c.Post(rel, headers, body)
	case http.MethodDelete:
		return c.Delete(rel, headers)
	default:
		return nil, errors.New("could not complete request: unsupported method")
	}
}

//Get performs a GET http request with the values being passed from the Do(). Since this request is only pulling
// information, it does not need to pass a body (body: nil).
func (c *Client) Get(rel *url.URL, headers http.Header) (*Response, error) {
	_url := c.BaseURL.ResolveReference(rel)
	request, err := http.NewRequest(http.MethodGet, _url.String(), nil)
	if err != nil {
		return nil, err
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

//Get performs a PUT http request with the values being passed from the Do(). Since this request is updating information,
// it passes a body (body: io.Reader).
func (c *Client) Put(rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	_url := c.BaseURL.ResolveReference(rel)
	request, err := http.NewRequest(http.MethodPut, _url.String(), body)
	if err != nil {
		return nil, err
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

//Get performs a POST http request with the values being passed from the Do(). Since this request is creating new
// information, it passes a body (body: io.Reader).
func (c *Client) Post(rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	_url := c.BaseURL.ResolveReference(rel)
	request, err := http.NewRequest(http.MethodPost, _url.String(), body)
	if err != nil {
		return nil, err
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

//Get performs a DELETE http request with the values being passed from the Do(). Since this request is removing
// information, it does not pass a body (body: nil).
func (c *Client) Delete(rel *url.URL, headers http.Header) (*Response, error) {
	_url := c.BaseURL.ResolveReference(rel)
	request, err := http.NewRequest(http.MethodDelete, _url.String(), nil)
	if err != nil {
		return nil, err
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

//do handles the retry logic for the request.
// Retry functionality is in place to retry a request that comes back with a server error response code if configured.
func (c *Client) do(request *http.Request) (*Response, error) {
	response, err := c.call(request)
	if err != nil {
		return nil, err
	}

	// attempt to retry a failed request if the the configuration for the max retry (MaxRetry int `json:"max-retry") is
	// greater than 0. Failed request only includes server response codes ranging between 500 - 599.
	if IsServerError(response.StatusCode) && c.MaxRetry > 0 {
		log.Warn("Received server response error code, retrying request...")
		retries := 0
		for retries < c.MaxRetry {
			time.Sleep(c.RetryDelay)

			retries++
			response, err = c.call(request)
			if err != nil {
				return nil, err
			}
		}
		log.Warn("Max retries reached")
	}
	return response, err
}

//call executes the actual request being made to the http client with the request that gets created (GET, PUT, POST, DELETE).
func (c *Client) call(request *http.Request) (*Response, error) {
	resp, err := c.HttpClient.Do(request)
	if err != nil {
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
func inRange(code, a, b int) bool {
	return a <= code && code < b
}

