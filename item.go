package ttlmap

import (
	"sync"
	"time"
)

// Item represents a record in the cache map
type Item struct {
	sync.RWMutex
	data    interface{}
	expires *time.Time
}

func (item *Item) touch(duration time.Duration) {
	item.Lock()
	defer item.Unlock()

	expiration := time.Now().Add(duration)
	item.expires = &expiration
}

func (item *Item) expired() bool {
	var value bool
	item.RLock()
	defer item.RUnlock()

	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	return value
}
