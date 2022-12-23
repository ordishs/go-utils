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
	expiryChan chan []V
}

// New creates a new ExpiringMap with the given expiry duration.

func New[K comparable, V any](expire time.Duration, expiryChan ...chan []V) *ExpiringMap[K, V] {
	m := &ExpiringMap[K, V]{
		expiry: expire,
		items:  make(map[K]*itemWrapper[V]),
	}

	if len(expiryChan) > 0 {
		m.expiryChan = expiryChan[0]
	}

	go func() {
		for {
			time.Sleep(expire)
			m.clean()
		}
	}()

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

	var expired []V

	for key, item := range m.items {
		if now > item.expiry {
			if m.expiryChan != nil {
				expired = append(expired, item.item)
			}

			delete(m.items, key)
		}
	}
	if m.expiryChan != nil && len(expired) > 0 {
		m.expiryChan <- expired
	}
}
