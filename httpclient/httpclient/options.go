package httpclient

import (
	"net/http"
	"time"
)

// Options represents the http client options
type Options struct {
	Retrier    Retriable
	Timeout    time.Duration
	RetryCount int
	Middles    []Middleware
	Transport  http.RoundTripper
}

func (opts Options) Clone() Options {
	opts1 := Options{
		Retrier:    opts.Retrier,
		Timeout:    opts.Timeout,
		RetryCount: opts.RetryCount,
	}

	for i := range opts.Middles {
		opts1.Middles = append(opts1.Middles, opts.Middles[i])
	}
	return opts1
}

// Option ...
type Option func(opts *Options)
