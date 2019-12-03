package options

import (
	"time"

	"github.com/DustyRat/post-it/pkg/client"
)

// Options ...
type Options struct {
	Input     string
	Output    string
	Latencies bool
	Flags     Flags

	Connections int
	Client      client.Config

	Headers     []string
	RawUrl      string
	RequestBody string

	Timeout            time.Duration
	IdleTimeout        time.Duration
	InsecureSkipVerify bool
}

// Flags ...
type Flags struct {
	Type    string
	Status  string
	Headers bool
	Body    bool
}
