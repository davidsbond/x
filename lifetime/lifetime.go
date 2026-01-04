// Package lifetime provides a mechanism for wrapping io.Closer implementations in a lifetime, automatically calling
// close after a specified duration.
package lifetime

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"
)

type (
	// The Lifetime type wraps an io.Closer implementation, making it accessible via the Value method until its specified
	// lifetime expires. After expiry, calls to Value will return ErrExpired.
	Lifetime[T io.Closer] struct {
		mutex      sync.RWMutex
		value      T
		expired    bool
		err        error
		cancel     context.CancelFunc
		expireOnce sync.Once
	}
)

var (
	// ErrExpired is the error given when calling Lifetime.Value on an expired lifetime.
	ErrExpired = errors.New("expired")
)

// New returns a new instance of the Lifetime type that wraps the given io.Closer. The specified lifetime is used to
// determine how long to wait until Close is called. After the lifetime has expired, calls to Lifetime.Value will
// return ErrExpired. If an error occurs calling Close it will be returned on the next call to Lifetime.Value.
func New[T io.Closer](value T, lifetime time.Duration) *Lifetime[T] {
	ctx, cancel := context.WithCancel(context.Background())

	lt := Lifetime[T]{
		value:  value,
		cancel: cancel,
	}

	go lt.wait(ctx, lifetime)

	return &lt
}

// Value returns T if it has not expired. Otherwise, it returns ErrExpired. If the lifetime has expired but calling
// Close failed, this method will return both ErrExpired and the error returned by Close as a joined error. Note that
// values you wrap may still be usable outside the scope of this lifetime, so you will still need to handle relevant
// close errors that may occur using the returned value in that scope.
func (lt *Lifetime[T]) Value() (T, error) {
	lt.mutex.RLock()
	defer lt.mutex.RUnlock()

	var zero T
	if lt.expired {
		return zero, errors.Join(ErrExpired, lt.err)
	}

	return lt.value, nil
}

// Expire causes immediate expiration of the underlying io.Closer.
func (lt *Lifetime[T]) Expire() {
	lt.expireOnce.Do(func() {
		lt.mutex.Lock()
		defer lt.mutex.Unlock()

		lt.expired = true
		lt.cancel()
		if err := lt.value.Close(); err != nil {
			lt.err = err
		}
	})
}

func (lt *Lifetime[T]) wait(ctx context.Context, lifetime time.Duration) {
	timer := time.NewTimer(lifetime)
	defer timer.Stop()

	// Here, we wait for either the lifetime to pass or if manual expiration has
	// been triggered. Once either of these conditions are met we close the io.Closer.
	select {
	case <-ctx.Done():
	case <-timer.C:
	}

	lt.Expire()
}
