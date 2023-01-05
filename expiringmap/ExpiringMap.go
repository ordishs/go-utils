package expiringmap

import (
	"sync"
	"time"
)

type itemWrapper[V any] struct {
	item   V
	expiry int64
}

// ExpiringMap is a map that expires items after a given duration.
// It uses go generics to allow for any type of key and value, although
// the key must be comparable.
// The map is safe for concurrent use.
type ExpiringMap[K comparable, V any] struct {
	mu         sync.RWMutex
	expiry     time.Duration
	items      map[K]*itemWrapper[V]
	evictionCh chan []V
	evictionFn func(V) bool
}

// New creates a new ExpiringMap with the given expiry duration.

func New[K comparable, V any](expire time.Duration) *ExpiringMap[K, V] {
	m := &ExpiringMap[K, V]{
		expiry: expire,
		items:  make(map[K]*itemWrapper[V]),
	}

	if expire != 0 {
		go func() {
			for {
				time.Sleep(expire)
				m.clean()
			}
		}()
	}

	return m
}

func (m *ExpiringMap[K, V]) WithEvictionChannel(ch chan []V) *ExpiringMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.evictionCh = ch
	return m
}

func (m *ExpiringMap[K, V]) WithEvictionFunction(f func(V) bool) *ExpiringMap[K, V] {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.evictionFn = f
	return m
}

// Set sets the value for the given key.
func (m *ExpiringMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = &itemWrapper[V]{
		item:   value,
		expiry: time.Now().Add(m.expiry).UnixNano(),
	}
}

// Get returns the value for the given key.
func (m *ExpiringMap[K, V]) Get(key K) (value V, ok bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	item, ok := m.items[key]
	if !ok {
		return
	}

	if time.Now().UnixNano() > item.expiry {
		ok = false
		return
	}

	value = item.item
	return
}

// Delete deletes the given key from the map.
func (m *ExpiringMap[K, V]) Delete(key K) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.items, key)
}

// Len returns the number of items in the map.
func (m *ExpiringMap[K, V]) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.items)
}

// Items returns a copy of the items in the map.
func (m *ExpiringMap[K, V]) Items() map[K]V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	items := make(map[K]V, len(m.items))

	for key, item := range m.items {
		items[key] = item.item
	}

	return items
}

// Clear clears the map.
func (m *ExpiringMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[K]*itemWrapper[V])
}

func (m *ExpiringMap[K, V]) clean() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixNano()

	var expiredItems []V

	for key, item := range m.items {
		if now > item.expiry {
			if m.evictionFn != nil && !m.evictionFn(item.item) {
				continue
			}

			if m.evictionCh != nil {
				expiredItems = append(expiredItems, item.item)
			}

			delete(m.items, key)
		}
	}

	if m.evictionCh != nil && len(expiredItems) > 0 {
		m.evictionCh <- expiredItems
	}
}
