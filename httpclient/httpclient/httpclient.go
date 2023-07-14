package httpclient

import (
	"io"
	"net/http"
	"net/url"
)

// DoFunc http client do func wrapper
type DoFunc func(*http.Request, func(*http.Response) error) error

// Middleware http client middleware
type Middleware func(DoFunc) DoFunc

// Client http client interface
type Client interface {
	Get(url string, headers http.Header, opts ...RequestOption) (*http.Response, error)
	Post(url string, body io.Reader, headers http.Header, opts ...RequestOption) (*http.Response, error)
	PostForm(url string, val url.Values, headers http.Header, opts ...RequestOption) (*http.Response, error)
	Put(url string, body io.Reader, headers http.Header, opts ...RequestOption) (*http.Response, error)
	Patch(url string, body io.Reader, headers http.Header, opts ...RequestOption) (*http.Response, error)
	Delete(url string, headers http.Header, opts ...RequestOption) (*http.Response, error)
	Do(req *http.Request) (*http.Response, error)
	DoWithOptions(request *http.Request, opts *RequestOptions) (*http.Response, error)
	Options() Options
	Clone(opts ...Option) Client
}
