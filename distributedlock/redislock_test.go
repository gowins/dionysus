package distributedlock

import (
	"testing"
)

func Test_redisLock_Lock(t *testing.T) {
	rl := redisLock{}
	rl.Lock()
}
