// Package keymux provides a keyed implementations of sync.Mutex that can be used to share locks for common keys. This
// can be useful for implementing global mutexes based off of identifiers like user identifiers or database names.
package keymux

import (
	"sync"
)

type (
	// The Mutex type is a wrapper around sync.Mutex that provides locking semantics using shared keys. Locking and
	// unlocking can be performed using a comparable type.
	Mutex[T comparable] struct {
		global   sync.RWMutex
		children map[T]*sync.Mutex
	}
)

// New returns a new Mutex.
func New[T comparable]() *Mutex[T] {
	return &Mutex[T]{
		children: make(map[T]*sync.Mutex),
	}
}

// Lock locks m for the specified key. If the lock is already in use, the calling goroutine blocks until the mutex is
// available.
func (m *Mutex[T]) Lock(key T) {
	m.get(key).Lock()
}

func (m *Mutex[T]) Unlock(key T) {
	m.get(key).Unlock()
}

// TryLock tries to lock m for the specified key and reports whether it succeeded.
func (m *Mutex[T]) TryLock(key T) bool {
	return m.get(key).TryLock()
}

func (m *Mutex[T]) get(key T) *sync.Mutex {
	m.global.RLock()
	child, ok := m.children[key]
	m.global.RUnlock()

	if !ok {
		m.global.Lock()
		child = &sync.Mutex{}
		m.children[key] = child
		m.global.Unlock()
	}

	return child
}
