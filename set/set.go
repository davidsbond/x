// Package set provides a simple set implementation. Sets act like maps but only contain values.
package set

import (
	"iter"
	"maps"
	"slices"
)

type (
	// The Set type is used like a map that only contains values. Values must implement the comparable builtin
	// interface.
	Set[T comparable] struct {
		entries map[T]struct{}
	}
)

// New returns a new Set.
func New[T comparable]() *Set[T] {
	return &Set[T]{
		entries: make(map[T]struct{}),
	}
}

// Put a value into the Set.
func (s *Set[T]) Put(v T) {
	s.entries[v] = struct{}{}
}

// Remove a value from the Set.
func (s *Set[T]) Remove(v T) {
	delete(s.entries, v)
}

// Len returns the number of entries within the Set.
func (s *Set[T]) Len() int {
	return len(s.entries)
}

// Values returns all values within the Set.
func (s *Set[T]) Values() []T {
	return slices.Collect(maps.Keys(s.entries))
}

// Clear all entries from the Set.
func (s *Set[T]) Clear() {
	s.entries = make(map[T]struct{})
}

// Range over all values in the Set.
func (s *Set[T]) Range() iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s.entries {
			if !yield(v) {
				return
			}
		}
	}
}
