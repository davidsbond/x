// Package slicepool provides a parameterized implementation of sync.Pool for use with slices.
package slicepool

import (
	"sync"
)

type (
	// The Pool type wraps a sync.Pool to expose parameterized implementations of Get and Put.
	Pool[T any] struct {
		pool sync.Pool
	}
)

// New returns a new Pool that will provide slices of type T with the specified length.
func New[T any](length int) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				b := make([]T, length)
				return &b
			},
		},
	}
}

// Get a slice of T from the Pool.
func (p *Pool[T]) Get() *[]T {
	buf := p.pool.Get().(*[]T)

	return buf
}

// Put returns a slice of T to the Pool.
func (p *Pool[T]) Put(buf *[]T) {
	p.pool.Put(buf)
}
