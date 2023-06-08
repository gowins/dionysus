package distributedlock

import (
	"strings"
	"testing"
)

func TestGetLockValue(t *testing.T) {
	lockValue, err := getLockValue()
	if err != nil {
		t.Errorf("want get lock value error nil, get error %v", err)
		return
	}
	if !strings.Contains(lockValue, "goid") {
		t.Errorf("want lockValue contains goid, get %v", lockValue)
	}
}
