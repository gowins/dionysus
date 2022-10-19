package registry

import "errors"

var (
	// Not found error when GetService is called
	ErrNotFound = errors.New("not found")
	// Watcher stopped error when watcher is stopped
	ErrWatcherStopped = errors.New("watcher stopped")
)
