package httpclient

import (
	"context"
	"net/http"
	"time"
)

// Options represents the http client options
type Options struct {
	Retrier      Retriable
	Timeout      time.Duration
	RetryCount   int
	Middles      []Middleware
	Transport    http.RoundTripper
	TracerEnable bool
}

func (opts Options) Clone() Options {
	opts1 := Options{
		Retrier:      opts.Retrier,
		Timeout:      opts.Timeout,
		RetryCount:   opts.RetryCount,
		TracerEnable: opts.TracerEnable,
	}

	for i := range opts.Middles {
		opts1.Middles = append(opts1.Middles, opts.Middles[i])
	}
	return opts1
}

// Option ...
type Option func(opts *Options)

type RequestOptions struct {
	Retrier    Retriable
	RetryCount int
	Ctx        context.Context
}

type RequestOption func(opts *RequestOptions)
