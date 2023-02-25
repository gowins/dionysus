package httpclient

import (
	"net/http"
	"time"

	"github.com/gowins/dionysus/httpclient/httpclient"
)

// Options http client options
type Options = httpclient.Options

// Option represents the http client options
type Option = httpclient.Option

// WithHTTPTimeout sets hystrix timeout
func WithHTTPTimeout(timeout time.Duration) Option {
	return func(c *Options) {
		c.Timeout = timeout
	}
}

// WithRetryCount sets the retry count for the client
func WithRetryCount(retryCount int) Option {
	return func(c *Options) {
		c.RetryCount = retryCount
	}
}

// WithRetrier sets the strategy for retrying
func WithRetrier(retrier httpclient.Retriable) Option {
	return func(c *Options) {
		c.Retrier = retrier
	}
}

// WithMiddleware sets the strategy for retrying
func WithMiddleware(m Middleware) Option {
	return func(c *Options) {
		c.Middles = append(c.Middles, m)
	}
}

// WithTransport sets the Transport of client
func WithTransport(t http.RoundTripper) Option {
	return func(c *Options) {
		c.Transport = t
	}
}
