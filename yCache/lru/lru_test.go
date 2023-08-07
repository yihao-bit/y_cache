package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("1234"))
	if v, ok := lru.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestCache_RemoveOldest(t *testing.T) {
	lru := New(int64(24), nil)
	lru.Add("key1", String("1234"))
	lru.Add("key2", String("1234"))
	lru.Add("key3", String("1234"))
	lru.Get("key1")
	lru.Add("key4", String("1234"))
	if _, ok := lru.Get("key2"); ok || lru.Len() != 3 {
		t.Fatalf("cache RemoveOldest failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	onEvicted := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64((10)), onEvicted)
	lru.Add("key1", String("1234"))
	lru.Add("key2", String("1234"))
	expect := []string{"key1"}
	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("not call onEvicted")
	}
}
