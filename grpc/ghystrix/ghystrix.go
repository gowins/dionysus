package ghystrix

import "github.com/afex/hystrix-go/hystrix"

const (
	defaultTimeout       = 10000
	defaultMaxConcurrent = 5000
)

func HystrixDefault(cfg hystrix.CommandConfig) {
	hystrix.DefaultTimeout = defaultTimeout
	hystrix.DefaultMaxConcurrent = defaultMaxConcurrent
	if cfg.Timeout != 0 {
		hystrix.DefaultTimeout = cfg.Timeout
	}
	if cfg.MaxConcurrentRequests != 0 {
		hystrix.DefaultMaxConcurrent = cfg.MaxConcurrentRequests
	}
	if cfg.ErrorPercentThreshold != 0 {
		hystrix.DefaultErrorPercentThreshold = cfg.ErrorPercentThreshold
	}
	if cfg.SleepWindow != 0 {
		hystrix.DefaultSleepWindow = cfg.SleepWindow
	}
	if cfg.RequestVolumeThreshold != 0 {
		hystrix.DefaultVolumeThreshold = cfg.RequestVolumeThreshold
	}
}

type HystrixCfg struct {
	Name string
	Cfg  hystrix.CommandConfig
}
