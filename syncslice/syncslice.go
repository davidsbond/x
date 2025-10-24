// Package syncslice provides a concurrent slice implementation using parameterized types.
package syncslice

import (
	"iter"
	"slices"
	"sync"
)

type (
	// The Slice type wraps a slice within a sync.RWMutex.
	Slice[T any] struct {
		mux   sync.RWMutex
		elems []T
	}
)

// New returns a new Slice with a specified length.
func New[T any](length uint) *Slice[T] {
	if length == 0 {
		return &Slice[T]{}
	}

	return &Slice[T]{
		elems: make([]T, length),
	}
}

// Len returns the number of elements within the Slice.
func (s *Slice[T]) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return len(s.elems)
}

// Append an element to the Slice.
func (s *Slice[T]) Append(v ...T) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.elems = append(s.elems, v...)
}

// At returns the element at the specified index. Panics if out-of-range like a regular slice would.
func (s *Slice[T]) At(idx uint) T {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.elems[idx]
}

// Set the element at a specified index to a given value. Panics if out-of-range like a regular slice would.
func (s *Slice[T]) Set(idx uint, v T) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.elems[idx] = v
}

// Range over all elements in the Slice.
func (s *Slice[T]) Range() iter.Seq2[uint, T] {
	return func(yield func(uint, T) bool) {
		s.mux.RLock()
		defer s.mux.RUnlock()

		for i, elem := range s.elems {
			if !yield(uint(i), elem) {
				return
			}
		}
	}
}

// Unwrap the Slice, returning a copy of the underlying slice.
func (s *Slice[T]) Unwrap() []T {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return slices.Clone(s.elems)
}
