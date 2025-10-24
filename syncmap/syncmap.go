// Package syncmap provides a concurrent map implementation using parameterized types.
package syncmap

import (
	"iter"
	"maps"
	"slices"
	"sync"
)

type (
	// The Map type wraps a map within a sync.RWMutex.
	Map[K comparable, V any] struct {
		mux     sync.RWMutex
		entries map[K]V
	}
)

// New returns a new Map.
func New[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{
		entries: make(map[K]V),
	}
}

// Put a value into the Map.
func (m *Map[K, V]) Put(k K, v V) {
	m.mux.Lock()
	m.entries[k] = v
	m.mux.Unlock()
}

// Get a value from the Map. The boolean return value indicates if a value was found.
func (m *Map[K, V]) Get(k K) (V, bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()

	v, ok := m.entries[k]
	return v, ok
}

// Values returns a slice of all values within the Map.
func (m *Map[K, V]) Values() []V {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return slices.Collect(maps.Values(m.entries))
}

// Keys returns a slice of all keys within the Map.
func (m *Map[K, V]) Keys() []K {
	m.mux.RLock()
	defer m.mux.RUnlock()

	return slices.Collect(maps.Keys(m.entries))
}

// Remove an entry from the Map.
func (m *Map[K, V]) Remove(k K) {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.entries, k)
}

// Len returns the number of entries within the Map.
func (m *Map[K, V]) Len() int {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return len(m.entries)
}

// Clear all entries within the Map.
func (m *Map[K, V]) Clear() {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.entries = make(map[K]V)
}

// Range over all keys and values in the Map.
func (m *Map[K, V]) Range() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.mux.RLock()
		defer m.mux.RUnlock()

		for k, v := range m.entries {
			if !yield(k, v) {
				return
			}
		}
	}
}
