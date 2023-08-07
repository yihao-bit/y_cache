package lru

import "container/list"

// cache 不是并发安全的
type Cache struct {
	maxBytes int64
	nbytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	//某条记录被移除时的回调函数
	onEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		nbytes:    0,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	e := c.ll.Back()
	if e != nil {
		c.ll.Remove(e)
		kv := e.Value
		en := kv.(*entry)
		delete(c.cache, en.key)
		c.nbytes -= int64(len(en.key)) + int64(en.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(en.key, en.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		oldLenV := kv.value.Len()
		kv.value = value
		c.nbytes += int64(value.Len()) - int64(oldLenV)
	} else {
		v := &entry{key: key, value: value}
		ele := c.ll.PushFront(v)
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
