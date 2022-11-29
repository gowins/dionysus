package memcache

import (
	"time"

	"github.com/allegro/bigcache/v3"
)

// ConfigOpt set bigcache.Config fileds
type ConfigOpt func(*bigcache.Config)

// WithShards set bigcache.Config.Shards field
func WithShards(shards int) ConfigOpt {
	return func(c *bigcache.Config) {
		c.Shards = shards
	}
}

// WithLifeWindow set bigcache.Config.LifeWindow field
func WithLifeWindow(lifeWindow time.Duration) ConfigOpt {
	return func(c *bigcache.Config) {
		c.LifeWindow = lifeWindow
	}
}

// WithLifeWindow set bigcache.Config.CleanWindow field
func WithCleanWindow(cleanWindow time.Duration) ConfigOpt {
	return func(c *bigcache.Config) {
		c.CleanWindow = cleanWindow
	}
}

// WithMaxEntriesInWindow set bigcache.Config.MaxEntriesInWindow field
func WithMaxEntriesInWindow(max int) ConfigOpt {
	return func(c *bigcache.Config) {
		c.MaxEntriesInWindow = max
	}
}

// WithMaxEntrySize set bigcache.Config.MaxEntrySize
func WithMaxEntrySize(max int) ConfigOpt {
	return func(c *bigcache.Config) {
		c.MaxEntrySize = max
	}
}

// WithStatsEnabled set bigcache.Config.StatsEnabled field
func WithStatsEnabled(enable bool) ConfigOpt {
	return func(c *bigcache.Config) {
		c.StatsEnabled = enable
	}
}

// WithVerbose set bigcache.Config.Verbose field
func WithVerbose(v bool) ConfigOpt {
	return func(c *bigcache.Config) {
		c.Verbose = v
	}
}

// WithHasher set bigcache.Config.Hasher field
func WithHasher(hasher bigcache.Hasher) ConfigOpt {
	return func(c *bigcache.Config) {
		c.Hasher = hasher
	}
}

// WithHardMaxCacheSize set bigcache.Config.HardMaxCacheSize field
func WithHardMaxCacheSize(size int) ConfigOpt {
	return func(c *bigcache.Config) {
		c.HardMaxCacheSize = size
	}
}

// WithOnRemove set bigcache.Config.OnRemove
func WithOnRemove(f func(string, []byte)) ConfigOpt {
	return func(c *bigcache.Config) {
		c.OnRemove = f
	}
}

// WithOnRemoveWithMetadata set bigcache.Config.OnRemoveWithMetadata field
func WithOnRemoveWithMetadata(f func(string, []byte, bigcache.Metadata)) ConfigOpt {
	return func(c *bigcache.Config) {
		c.OnRemoveWithMetadata = f
	}
}

// WithOnRemoveWithReason set bigcache.Config.OnRemoveWithReason field
func WithOnRemoveWithReason(f func(string, []byte, bigcache.RemoveReason)) ConfigOpt {
	return func(c *bigcache.Config) {
		c.OnRemoveWithReason = f
	}
}

// WithLogger set bigcache.Config.Logger field
func WithLogger(l bigcache.Logger) ConfigOpt {
	return func(c *bigcache.Config) {
		c.Logger = l
	}
}
