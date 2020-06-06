package options

import (
	"time"

	"github.com/DustyRat/post-it/internal/http"
)

// Options ...
type Options struct {
	Input     string
	Output    string
	Histogram bool
	Latency   bool
	Flags     Flags

	Connections int
	Client      http.Config

	Headers     []string
	RawUrl      string
	RequestBody string

	Timeout            time.Duration
	IdleTimeout        time.Duration
	InsecureSkipVerify bool
}

// Flags ...
type Flags struct {
	Status  string
	Errors  bool
	Headers bool
	Body    bool
}
