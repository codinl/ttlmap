package ttlmap

import (
	"sync"
	"time"
)

// TtlMap is a synchronised map of items that auto-expire once stale
type TtlMap struct {
	mutex sync.RWMutex
	ttl   time.Duration
	items map[interface{}]*Item
}

// Set is a thread-safe way to add new items to the map
func (t *TtlMap) Set(key interface{}, data interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	item := &Item{data: data}
	item.touch(t.ttl)
	t.items[key] = item
}

// Get is a thread-safe way to lookup items
// Every lookup, also touches the item, hence extending it's life
func (t *TtlMap) Get(key interface{}) (data interface{}, found bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	item, exists := t.items[key]
	if !exists || item.expired() {
		data = ""
		found = false
	} else {
		item.touch(t.ttl)
		data = item.data
		found = true
	}
	return
}

// Count returns the number of items in the cache
// (helpful for tracking memory leaks)
func (t *TtlMap) Count() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	count := len(t.items)
	return count
}

func (t *TtlMap) cleanup() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for key, item := range t.items {
		if item.expired() {
			delete(t.items, key)
		}
	}
}

func (t *TtlMap) startCleanupTimer() {
	duration := t.ttl
	if duration < time.Second {
		duration = time.Second
	}
	ticker := time.Tick(duration)
	go (func() {
		for {
			select {
			case <-ticker:
				t.cleanup()
			}
		}
	})()
}

// NewTtlMap is a helper to create instance of the Map struct
func NewTtlMap(duration time.Duration) *TtlMap {
	t := &TtlMap{
		ttl:   duration,
		items: map[interface{}]*Item{},
	}
	t.startCleanupTimer()
	return t
}
