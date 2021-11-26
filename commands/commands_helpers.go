package commands

import (
	"sync"
	"time"

	"github.com/arraybot/nautilus/database"
)

// ttlItem wraps the cache value with a creation time.
type ttlItem struct {
	v bool
	c int64
}

// ttlCache is a cache where entities are only valid for a certain time after creation.
type ttlCache struct {
	m map[string]*ttlItem
	f func(string) bool
	i int64
	l sync.Mutex
}

// get gets a value from the cache, fetching it if it is not present or expired.
func (t *ttlCache) get(k string) bool {
	// Get the current time in milliseconds.
	n := time.Now().UnixNano() / int64(time.Millisecond)
	t.l.Lock()
	// Get the entry.
	v, c := t.m[k]
	// If it doesn't exist, or too much time since creation has elapsed.
	if !c || (n-v.c) > t.i {
		// Create a new item by fetching it again.
		v = &ttlItem{
			v: t.f(k),
			c: n,
		}
	}
	// Update the value and return.
	t.m[k] = v
	t.l.Unlock()
	return v.v
}

// invalidate invalidates the cache entry for a given string.
func (t *ttlCache) invalidate(k string) {
	t.l.Lock()
	delete(t.m, k)
	t.l.Unlock()
}

// Caches whether or not to send commands invsible to the server.
// Valid for: 10 minutes.
var cacheInvisibility = ttlCache{
	m: make(map[string]*ttlItem),
	f: func(k string) bool {
		return database.ReplyHidden(k)
	},
	i: 600 * 1000,
}
