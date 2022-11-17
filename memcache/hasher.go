package memcache

import "github.com/cespare/xxhash/v2"

type hasher struct{}

func (h *hasher) Sum64(key string) uint64 {
	return xxhash.Sum64String(key)
}
