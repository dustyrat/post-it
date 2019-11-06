package options

import (
	"time"

	"github.com/DustyRat/post-it/pkg/client"
)

type Options struct {
	Input  string
	Output string
	Flags  Flags

	Connections int
	Client      client.Config

	Headers     []string
	RawUrl      string
	RequestBody string

	Timeout            time.Duration
	IdleTimeout        time.Duration
	InsecureSkipVerify bool
}

type Flags struct {
	Type    string
	Status  string
	Headers bool
	Body    bool
}
