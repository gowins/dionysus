package httpclient

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

func createTransport() http.RoundTripper {
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          runtime.GOMAXPROCS(0)*12 + 100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0)*12 + 1, // http client缓存数量
	}

	return trans
}
