package singleflight

import "sync"

type Call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*Call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*Call)
	}
	if v, ok := g.m[key]; ok {
		g.mu.Unlock()
		//等待正在进行的请求结束，获取这个请求的返回值
		v.wg.Wait()
		return v.val, v.err
	}
	c := new(Call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return c.val, c.err
}
