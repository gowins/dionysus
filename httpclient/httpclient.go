package httpclient

import (
	"bytes"
	"context"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/gowins/dionysus/httpclient/httpclient"
)

const (
	defaultRetryCount  = 1
	defaultHTTPTimeout = 2 * time.Second
)

type DoFunc = httpclient.DoFunc
type Middleware = httpclient.Middleware
type Client = httpclient.Client

var _ Client = (*client)(nil)

// client is the Client implementation
type client struct {
	client *http.Client
	opts   Options
	do     DoFunc
}

func (c *client) getRequestOptions(opts ...RequestOption) *RequestOptions {
	reqOptions := &RequestOptions{
		Retrier:    c.opts.Retrier,
		RetryCount: c.opts.RetryCount,
		Ctx:        context.Background(),
	}
	for _, opt := range opts {
		opt(reqOptions)
	}
	return reqOptions
}

func (c *client) PostForm(url string, val url.Values, headers http.Header, opts ...RequestOption) (*http.Response, error) {
	requestOptions := c.getRequestOptions(opts...)
	request, err := http.NewRequestWithContext(requestOptions.Ctx, http.MethodPost, url, strings.NewReader(val.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "PostForm - request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return c.DoWithOptions(request, requestOptions)
}

func (c *client) Clone(opts ...httpclient.Option) Client {
	opts1 := c.opts.Clone()
	for _, opt := range opts {
		opt(&opts1)
	}

	return newClient(opts1)
}

func (c *client) Options() Options {
	return c.opts
}

func newOptions(opts ...Option) Options {
	opts1 := Options{
		Timeout:    defaultHTTPTimeout,
		RetryCount: defaultRetryCount,
		Retrier:    NewRetrier(NewConstantBackoff(10*time.Millisecond, 50*time.Millisecond)),
		Transport:  createTransport(),
	}

	for _, opt := range opts {
		opt(&opts1)
	}

	if opts1.TracerEnable {
		opts1.Transport = otelhttp.NewTransport(opts1.Transport)
	}

	return opts1
}

// New returns a new instance of http client
func New(opts ...Option) Client {
	opts1 := newOptions(opts...)
	return newClient(opts1)
}

// NewWithTracer returns a new instance of http client with tracer
func NewWithTracer(opts ...Option) Client {
	opts = append(opts, WithTracerEnable())
	return New(opts...)
}

func newClient(opts Options) Client {
	c := &client{
		opts: opts,
		client: &http.Client{
			Timeout:   opts.Timeout,
			Transport: opts.Transport,
		},
	}

	do := c.doFunc
	for i := len(opts.Middles); i > 0; i-- {
		do = opts.Middles[i-1](do)
	}
	c.do = do

	return c
}

// Get makes a HTTP GET request to provided URL
func (c *client) Get(url string, headers http.Header, opts ...RequestOption) (*http.Response, error) {
	requestOptions := c.getRequestOptions(opts...)
	request, err := http.NewRequestWithContext(requestOptions.Ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "GET - request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	return c.DoWithOptions(request, requestOptions)
}

// Post makes a HTTP POST request to provided URL and requestBody
func (c *client) Post(url string, body io.Reader, headers http.Header, opts ...RequestOption) (*http.Response, error) {
	requestOptions := c.getRequestOptions(opts...)
	request, err := http.NewRequestWithContext(requestOptions.Ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "POST - request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	return c.DoWithOptions(request, requestOptions)
}

// Put makes a HTTP PUT request to provided URL and requestBody
func (c *client) Put(url string, body io.Reader, headers http.Header, opts ...RequestOption) (*http.Response, error) {
	requestOptions := c.getRequestOptions(opts...)
	request, err := http.NewRequestWithContext(requestOptions.Ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PUT - request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	return c.DoWithOptions(request, requestOptions)
}

// Patch makes a HTTP PATCH request to provided URL and requestBody
func (c *client) Patch(url string, body io.Reader, headers http.Header, opts ...RequestOption) (*http.Response, error) {
	requestOptions := c.getRequestOptions(opts...)
	request, err := http.NewRequestWithContext(requestOptions.Ctx, http.MethodPatch, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "PATCH - request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	return c.DoWithOptions(request, requestOptions)
}

// Delete makes a HTTP DELETE request with provided URL
func (c *client) Delete(url string, headers http.Header, opts ...RequestOption) (*http.Response, error) {
	requestOptions := c.getRequestOptions(opts...)
	request, err := http.NewRequestWithContext(requestOptions.Ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "DELETE - request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}

	return c.DoWithOptions(request, requestOptions)
}

func (c *client) doFunc(req *http.Request, fn func(*http.Response) error) error {
	// nolint:bodyclose
	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return fn(response)
}

// Do makes an HTTP request with the native `http.Do` interface
func (c *client) Do(request *http.Request) (*http.Response, error) {
	opts := &RequestOptions{
		Retrier:    c.opts.Retrier,
		RetryCount: c.opts.RetryCount,
	}
	return c.DoWithOptions(request, opts)
}

// DoWithOptions makes an HTTP request with the native `http.DoWithOptions` interface
func (c *client) DoWithOptions(request *http.Request, opts *RequestOptions) (*http.Response, error) {
	var (
		err        error
		resp       *http.Response
		bodyReader *bytes.Reader
		reqData    []byte
	)

	if opts == nil {
		opts = &RequestOptions{
			Retrier:    c.opts.Retrier,
			RetryCount: c.opts.RetryCount,
		}
	}

	if request.Body != nil {
		reqData, err = io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(reqData)
		request.Body = io.NopCloser(bodyReader) // prevents closing the body between retries
	}

	for i := 0; i < opts.RetryCount; i++ {
		if resp != nil {
			resp.Body.Close()
		}

		err = c.do(request, func(response *http.Response) error {
			resp = response
			return nil
		})

		if err != nil {
			if backoffTime := opts.Retrier.NextInterval(i); backoffTime != 0 {
				time.Sleep(backoffTime)
			}
			if len(reqData) != 0 {
				bodyReader = bytes.NewReader(reqData)
				request.Body = io.NopCloser(bodyReader)
			}
			continue
		}

		break
	}

	return resp, err
}
