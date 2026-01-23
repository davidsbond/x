// Package closer provides types and functions for working with io.Closer implementations. Applications typically end
// up with a lot of io.Closer implementations in their boot and shutdown sequences, these utilities aim to make that
// a bit less cumbersome.
package closer

import (
	"errors"
	"io"
	"slices"
)

type (
	// The Collection type is responsible for managing multiple io.Closer implementations.
	Collection struct {
		closers []io.Closer
	}
)

// NewCollection returns a new instance of the Collection type prepopulated with the provided io.Closer implementations.
func NewCollection(closers ...io.Closer) *Collection {
	return &Collection{
		closers: closers,
	}
}

// Add an io.Closer implementation to the Collection.
func (co *Collection) Add(c io.Closer) {
	co.closers = append(co.closers, c)
}

// Close all stored io.Closer implementations in reverse order of their registration. Reverse order is desired. For
// example, when running an HTTP API you'll likely create your database connection first before your HTTP server. When
// performing a graceful shutdown, you'll want to stop your HTTP server before you stop your database connection.
// Otherwise, requests still in progress would not be able to communicate with the database.
func (co *Collection) Close() error {
	errs := make([]error, len(co.closers))
	for i, c := range slices.Backward(co.closers) {
		errs[i] = c.Close()
	}

	return errors.Join(errs...)
}
