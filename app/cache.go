package app

import (
	"killl/lib/log"
	"sync"
)

type Cache struct {
	cache map[string]struct{}
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]struct{}),
	}
}

func (c *Cache) put(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[key] = struct{}{}
}

func (c *Cache) del(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.cache, key)
}

func (c *Cache) exist(key string) bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if _, ok := c.cache[key]; !ok {
		return false
	}
	return true
}

func (c *Cache) clear(now map[string]struct{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	temp := make([]string, 0)
	for key := range c.cache {
		if _, ok := now[key]; !ok {
			temp = append(temp, key)
		}
	}
	for _, name := range temp {
		delete(c.cache, name)
		log.Debugf("cache clear: id(%s)", name)
	}
}
