package memcache

import "github.com/cespare/xxhash/v2"

type Hasher struct{}

func (h *Hasher) Sum64(key string) uint64 {
	return xxhash.Sum64String(key)
}
