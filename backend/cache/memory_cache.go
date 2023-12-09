package cache

import (
	"errors"
	"main/config"
	"sync"
	"time"
)

type MemoryCache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]MemoryCacheItem
}

type MemoryCacheItem struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	items := make(map[string]MemoryCacheItem)
	cache := MemoryCache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

func (c *MemoryCache) Set(key string, value interface{}, duration time.Duration) {
	var expiration int64

	if duration == 0 {
		duration = c.defaultExpiration
	}

	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()

	defer c.Unlock()

	c.items[key] = MemoryCacheItem{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

func (c *MemoryCache) Delete(key string) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

func (c *MemoryCache) StartGC() {
	go c.GC()
}

func (c *MemoryCache) GC() {
	for {
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

func (c *MemoryCache) expiredKeys() (keys []string) {
	c.RLock()
	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

func (c *MemoryCache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}

var MemoryCacheInstance = NewMemoryCache(
	time.Duration(config.Settings.MemoryCacheDurationMin)*time.Minute,
	time.Duration(config.Settings.MemoryCacheDurationMin+5)*time.Minute,
)
