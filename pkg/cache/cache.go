package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	*gocache.Cache
}

func New(defaultExpiration time.Duration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		Cache: gocache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *Cache) SetOrGet(key string, duration time.Duration, setFunc func() interface{}) interface{} {
	valFromCache, incache := c.Cache.Get(key)
	if !incache {
		valToCache := setFunc()
		c.Cache.Set(key, valToCache, duration)
		return valToCache
	}

	return valFromCache
}
