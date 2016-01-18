package ttlmap

import (
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	cache := &TtlMap{
		ttl:   time.Second,
		items: map[interface{}]*Item{},
	}

	data, exists := cache.Get("hello")
	if exists || data != "" {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world")
	data, exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}

	type M struct {
		Name string
	}

	model := M {
		Name:"name",
	}

	cache.Set(1,model)
	m, exists := cache.Get(1)
	if !exists {
		t.Errorf("Expected cache to return data for `model`")
	}
	if m.(M).Name != "name" {
		t.Errorf("Expected cache m.Name == 'name'")
	}

}

func TestExpiration(t *testing.T) {
	cache := &TtlMap{
		ttl:   time.Second,
		items: map[interface{}]*Item{},
	}

	cache.Set("x", 1)
	cache.Set("y", "z")
	cache.Set("z", 1.0)
	cache.startCleanupTimer()

	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected cache to contain 3 items")
	}

	<-time.After(500 * time.Millisecond)
	cache.mutex.Lock()
	cache.items["y"].touch(time.Second)
	item, exists := cache.items["x"]
	cache.mutex.Unlock()
	if !exists || item.data != 1 || item.expired() {
		t.Errorf("Expected `x` to not have expired after 200ms")
	}

	<-time.After(time.Second)
	cache.mutex.RLock()
	_, exists = cache.items["x"]
	if exists {
		t.Errorf("Expected `x` to have expired")
	}
	_, exists = cache.items["z"]
	if exists {
		t.Errorf("Expected `z` to have expired")
	}
	_, exists = cache.items["y"]
	if !exists {
		t.Errorf("Expected `y` to not have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 1 {
		t.Errorf("Expected cache to contain 1 item")
	}

	<-time.After(600 * time.Millisecond)
	cache.mutex.RLock()
	_, exists = cache.items["y"]
	if exists {
		t.Errorf("Expected `y` to have expired")
	}
	cache.mutex.RUnlock()

	count = cache.Count()
	if count != 0 {
		t.Errorf("Expected cache to be empty")
	}
}
