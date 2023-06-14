package distributedlock

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

var goroutineSpace = []byte("goroutine ")

func curGoroutineID() (string, error) {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	// Parse the 4707 out of "goroutine 4707 ["
	b = bytes.TrimPrefix(b, goroutineSpace)
	i := bytes.IndexByte(b, ' ')
	if i < 0 {
		return "", fmt.Errorf("no space found in %q", b)
	}
	b = b[:i]
	return string(b), nil
}

func GetLockValue() (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("get host name error %v", err)
	}
	goid, err := curGoroutineID()
	if err != nil {
		return "", fmt.Errorf("get go id error %v", err)
	}
	return hostName + "goid" + goid, nil
}
