package yCache

import (
	"fmt"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGetGroup(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	group := NewGroup("test", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))
	for k, v := range db {
		if view, err := group.Get(k); err != nil || string(view.b) != v {
			t.Fatalf("failed get key of tom or jack or sam")
		}
		if _, err := group.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache missed")
		}
	}
	if view, err := group.Get("unkonwn"); view.b != nil || err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
