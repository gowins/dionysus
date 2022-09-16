package rmq

import "github.com/apache/rocketmq-client-go/v2/rlog"

func SetLogger(logger rlog.Logger) {
	rlog.SetLogLevel("info")
	rlog.SetLogger(logger)
}
