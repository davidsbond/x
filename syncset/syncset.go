// Package syncset provides a concurrent wrapper around the set.Set type.
package syncset

import (
	"iter"
	"sync"

	"github.com/davidsbond/x/set"
)

type (
	// The Set type wraps a set.Set within a sync.RWMutex.
	Set[T comparable] struct {
		mux sync.RWMutex
		s   *set.Set[T]
	}
)

// New returns a new Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		s: set.New[T](),
	}
}

// Put a value into the Set.
func (s *Set[T]) Put(v T) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.s.Put(v)
}

// Remove a value from the Set.
func (s *Set[T]) Remove(v T) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s.Remove(v)
}

// Len returns the number of entries within the Set.
func (s *Set[T]) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.s.Len()
}

// Values returns all values within the Set.
func (s *Set[T]) Values() []T {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.s.Values()
}

// Clear all entries from the Set.
func (s *Set[T]) Clear() {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s.Clear()
}

// Contains returns true if a given value is present within the Set.
func (s *Set[T]) Contains(v T) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.s.Contains(v)
}

// Range over all values in the Set.
func (s *Set[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.mux.RLock()
		defer s.mux.RUnlock()

		for v := range s.s.Range() {
			if !yield(v) {
				return
			}
		}
	}
}
