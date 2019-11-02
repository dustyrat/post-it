package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	Method   string
	Header   http.Header
	URL      *url.URL
	Body     io.Reader
	Response *Response
}

func NewRequest(method, rawurl string, header http.Header, body io.Reader, fields map[string]string) (*Request, error) {
	for k, v := range fields {
		rawurl = strings.Replace(rawurl, fmt.Sprintf("{%s}", k), v, 1)
	}

	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	return &Request{
		Method: method,
		Header: header,
		URL:    uri,
		Body:   body,
	}, nil
}

func ParseHeaders(headers []string) http.Header {
	header := http.Header{}
	for _, h := range headers {
		head := strings.Split(h, ":")
		for _, v := range strings.Split(strings.TrimSpace(head[1]), ",") {
			header.Add(head[0], strings.TrimSpace(v))
		}
	}
	return header
}
