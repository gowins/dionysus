package ghystrix

import (
	"testing"
	"time"

	"github.com/afex/hystrix-go/hystrix"
)

func TestHystrixDefault(t *testing.T) {
	HystrixDefault(hystrix.CommandConfig{
		Timeout:                int(time.Second),
		MaxConcurrentRequests:  3,
		ErrorPercentThreshold:  1,
		SleepWindow:            5,
		RequestVolumeThreshold: 5,
	})
}
