package ginx

import (
	"github.com/gin-gonic/gin"
)

// GinOption set properties of gin.Engine instance
type GinOption func(*gin.Engine)

// WithRedirectTrailingSlash set RedirectTrailingSlash
func WithRedirectTrailingSlash(redirect bool) GinOption {
	return func(e *gin.Engine) {
		e.RedirectTrailingSlash = redirect
	}
}

// WithRedirectFixedPath set RedirectFixedPath
func WithRedirectFixedPath(redirect bool) GinOption {
	return func(e *gin.Engine) {
		e.RedirectFixedPath = redirect
	}
}

// WithHandleMethoNotAllowed set HandleMethodNotAllowed
func WithHandleMethoNotAllowed(allowed bool) GinOption {
	return func(e *gin.Engine) {
		e.HandleMethodNotAllowed = allowed
	}
}

// WithForwardedByClientIP set ForwardedByClientIP
func WithForwardedByClientIP(forwarded bool) GinOption {
	return func(e *gin.Engine) {
		e.ForwardedByClientIP = forwarded
	}
}

// WithUseRawPath set UseRawPath
func WithUseRawPath(rawPath bool) GinOption {
	return func(e *gin.Engine) {
		e.UseRawPath = rawPath
	}
}

// WithUnescapePathValues set UnescapePathValues
func WithUnescapePathValues(unescape bool) GinOption {
	return func(e *gin.Engine) {
		e.UnescapePathValues = unescape
	}
}

// WithRemoveExtraSlash set RemoveExtraSlash
func WithRemoveExtraSlash(remove bool) GinOption {
	return func(e *gin.Engine) {
		e.RemoveExtraSlash = remove
	}
}

// WithRemoteIPHeaders set RemoteIPHeaders
func WithRemoteIPHeaders(headers []string) GinOption {
	return func(e *gin.Engine) {
		e.RemoteIPHeaders = headers
	}
}

// WithTrustedPlatform set TrustedPlatform
func WithTrustedPlatform(trusted string) GinOption {
	return func(e *gin.Engine) {
		e.TrustedPlatform = trusted
	}
}

// WithMaxMultipartMemory set MaxMultipartMemory
func WithMaxMultipartMemory(memory int64) GinOption {
	return func(e *gin.Engine) {
		e.MaxMultipartMemory = memory
	}
}

// WithUseH2C set UseH2C
func WithUseH2C(h2c bool) GinOption {
	return func(e *gin.Engine) {
		e.UseH2C = h2c
	}
}

// WithContextWithFallback set ContextWithFallback
func WithContextWithFallback(fallback bool) GinOption {
	return func(e *gin.Engine) {
		e.ContextWithFallback = fallback
	}
}
