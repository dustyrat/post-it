package http

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

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
	Request          *http.Request
}

// Client ...
type Client struct {
	client  *http.Client
	url     *url.URL
	headers http.Header
}

var (
	registry                       = prometheus.NewRegistry()
	Registry prometheus.Registerer = registry
	Gatherer prometheus.Gatherer   = registry
	m        *metrics
)

type metrics struct {
	status   *prometheus.CounterVec
	duration *prometheus.HistogramVec
	summary  *prometheus.SummaryVec
}

func init() {
	m = &metrics{
		status: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_outbound_requests_status_total",
				Help: "Counter of successful Outbound HTTP requests by status code.",
			},
			[]string{"method", "code"},
		),
		duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_outbound_request_duration_seconds",
				Help: "Histogram of latencies for Outbound HTTP requests.",
				Buckets: []float64{
					.001, .0025, .005, .0075,
					.01, .025, .05, .075,
					.1, .25, .5, .75,
					1, 2.5, 5, 7.5,
					10,
				},
			},
			[]string{"method"},
		),
		summary: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "http_outbound_request_quantile",
				Help:       "Summary of latencies for Outbound HTTP requests.",
				Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001, 1.00: 0.0},
				MaxAge:     7 * 24 * time.Hour,
			},
			[]string{"method"},
		),
	}
	Registry.MustRegister(m.status, m.duration, m.summary)
}

// Config ...
type Config struct {
	URL string

	// Timeout specifies a time limit for requests made by this
	// Client. The timeout includes connection time, any
	// redirects, and reading the response body. The timer remains
	// running after Get, Head, Post, or Do return and will
	// interrupt reading of the Response.Body.
	// A Timeout of zero means no timeout.
	Timeout time.Duration

	// IdleConnTimeout is the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing
	// itself.
	// Zero means no limit.
	IdleConnTimeout time.Duration

	// InsecureSkipVerify controls whether a client verifies the
	// server's certificate chain and host name.
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	InsecureSkipVerify bool

	// MaxConnsPerHost optionally limits the total number of
	// connections per host, including connections in the dialing,
	// active, and idle states. On limit violation, dials will block.
	// Zero means no limit.
	MaxConnsPerHost int

	// MaxIdleConns controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	MaxIdleConns int

	// MaxIdleConnsPerHost, if non-zero, controls the maximum idle
	// (keep-alive) connections to keep per-host. If zero,
	// DefaultMaxIdleConnsPerHost is used.
	MaxIdleConnsPerHost int

	Headers http.Header
}

// New ...
func New(conf Config) (*Client, error) {
	uri, err := url.Parse(conf.URL)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: conf.InsecureSkipVerify,
		},
		MaxConnsPerHost:     conf.MaxConnsPerHost,
		MaxIdleConns:        conf.MaxIdleConns,
		MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		IdleConnTimeout:     conf.IdleConnTimeout * time.Millisecond,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   conf.Timeout * time.Millisecond,
	}

	return &Client{client: client, url: uri, headers: conf.Headers}, nil
}

// Do ...
func (c *Client) Do(method string, rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	uri := c.url.ResolveReference(rel)
	request, err := http.NewRequest(method, uri.String(), body)
	if err != nil {
		return nil, err
	}

	for k, vs := range c.headers {
		for _, v := range vs {
			headers.Add(k, v)
		}
	}
	request.Header = headers
	return c.do(request)
}

// Get helper method for making a GET request
func (c *Client) Get(rel *url.URL, headers http.Header) (*Response, error) {
	return c.Do(http.MethodGet, rel, headers, nil)
}

// Put helper method for making a PUT request
func (c *Client) Put(rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	return c.Do(http.MethodPut, rel, headers, body)
}

// Post helper method for making a POST request
func (c *Client) Post(rel *url.URL, headers http.Header, body io.Reader) (*Response, error) {
	return c.Do(http.MethodPost, rel, headers, body)
}

// Delete helper method for making a DELETE request
func (c *Client) Delete(rel *url.URL, headers http.Header) (*Response, error) {
	return c.Do(http.MethodDelete, rel, headers, nil)
}

// Head helper method for making a HEAD request
func (c *Client) Head(rel *url.URL, headers http.Header) (*Response, error) {
	return c.Do(http.MethodHead, rel, headers, nil)
}

func (c *Client) do(request *http.Request) (*Response, error) {
	start := time.Now()
	resp, err := c.client.Do(request)
	if err != nil {
		m.status.WithLabelValues(strings.ToLower(request.Method), "0").Inc()
		if err, ok := err.(*url.Error); ok {
			return nil, err.Unwrap()
		}
		return nil, err
	}
	defer resp.Body.Close()
	m.status.WithLabelValues(strings.ToLower(request.Method), strconv.Itoa(resp.StatusCode)).Inc()

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
		Request:          request,
	}
	m.duration.WithLabelValues(strings.ToLower(request.Method)).Observe(response.Duration.Seconds())
	m.summary.WithLabelValues(strings.ToLower(request.Method)).Observe(response.Duration.Seconds())
	return &response, nil
}

// IsSuccessful checks if server code being passed is a successfully code
func IsSuccessful(code int) bool {
	return inRange(code, 200, 300)
}

// IsClientError checks if server code being passed is a client error code
func IsClientError(code int) bool {
	return inRange(code, 400, 500)
}

// IsServerError checks if server code being passed is a server error code
func IsServerError(code int) bool {
	return inRange(code, 500, 600)
}

func inRange(code, a, b int) bool {
	return a <= code && code < b
}

// InRange ...
func InRange(code, a, b int) bool {
	return inRange(code, a, b)
}
