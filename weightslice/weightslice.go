// Package weightslice provides a "weighted slice" implementation where the order of elements within the slice changes
// based on a comparable weight value.
package weightslice

import (
	"cmp"
	"iter"
	"slices"
	"sync"
)

type (
	// The Slice type represents a slice of values that each have an associated weight that modifies their order within
	// the slice. The Slice is sorted on calls to Append and SetWeight.
	Slice[T any, W cmp.Ordered] struct {
		mutex    sync.RWMutex
		elements []entry[T, W]
		compare  func(a, b entry[T, W]) int
	}

	entry[V any, Weight comparable] struct {
		value  V
		weight Weight
	}
)

const (
	// Ascending can be passed to New to order elements in ascending order.
	Ascending uint = iota
	// Descending can be passed to New to order elements in descending order.
	Descending
)

// New returns a new Slice instance that will sort elements by weight in the specified sort order. The sort order can
// either be Ascending or Descending. For any other value, Ascending is assumed.
func New[T any, W cmp.Ordered](sort uint) *Slice[T, W] {
	compare := func(a, b entry[T, W]) int {
		return cmp.Compare(a.weight, b.weight)
	}

	if sort == Descending {
		compare = func(a, b entry[T, W]) int {
			return cmp.Compare(b.weight, a.weight)
		}
	}

	return &Slice[T, W]{
		elements: make([]entry[T, W], 0),
		compare:  compare,
	}
}

// Append an element to the Slice with an initial weight.
func (s *Slice[T, W]) Append(value T, weight W) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	e := entry[T, W]{
		value:  value,
		weight: weight,
	}

	s.elements = append(s.elements, e)
	slices.SortFunc(s.elements, s.compare)
}

// SetWeight modifies the weight value of an element at a specified index. This method will panic if the index is
// out of bounds.
func (s *Slice[T, W]) SetWeight(idx int, weight W) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.elements[idx].weight = weight
	slices.SortFunc(s.elements, s.compare)
}

// Range over all elements in the Slice based on their weight. This method uses a copy of the original slice for
// iteration to avoid concurrency issues.
func (s *Slice[T, W]) Range() iter.Seq2[int, T] {
	s.mutex.RLock()
	snapshot := slices.Clone(s.elements)
	s.mutex.RUnlock()

	return func(yield func(int, T) bool) {
		for i, e := range snapshot {
			if !yield(i, e.value) {
				return
			}
		}
	}
}
