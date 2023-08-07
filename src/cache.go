package src

import (
	"sync"
	"y_cache/src/lru"
)

type cache struct {
	//这里不能用读写锁，get方法因为淘汰策略有写的动作
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), true
	}
	return
}
