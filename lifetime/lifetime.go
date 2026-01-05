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
		reset      chan time.Duration
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
		// This channel is used to signal that the timer is to be reset, it must be kept as a buffered size of 1
		// and is used in a "last write wins" fashion. When calling Lifetime.Reset we'll drain this channel if
		// there are pending writes.
		reset: make(chan time.Duration, 1),
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

// Reset the lifetime to a new duration. This modifies the expiration to occur after the given duration and does not
// add additional time to any remaining lifetime.
func (lt *Lifetime[T]) Reset(lifetime time.Duration) error {
	lt.mutex.RLock()

	if lt.expired {
		lt.mutex.RUnlock()
		return ErrExpired
	}

	lt.mutex.RUnlock()

	// To keep this method from blocking, we'll attempt to write directly to the reset channel, if it blocks, we'll
	// drain it ourselves and rewrite to it.
	select {
	case lt.reset <- lifetime:
		break
	default:
		<-lt.reset
		lt.reset <- lifetime
	}

	return nil
}

func (lt *Lifetime[T]) wait(ctx context.Context, lifetime time.Duration) {
	timer := time.NewTimer(lifetime)
	defer timer.Stop()

	for {
		// The two conditions for lifetimes expiring is the context is canceled (A call to Lifetime.Expire) or the
		// lifetime has expired naturally.
		select {
		case <-ctx.Done():
			lt.Expire()
			return
		case <-timer.C:
			lt.Expire()
			return
		case lifetime = <-lt.reset:
			if !timer.Stop() {
				select {
				// If we can't stop the timer because it has a queued tick in its channel, we'll first drain that
				// then call Reset.
				case <-timer.C:
					break
				default:
					break
				}
			}

			timer.Reset(lifetime)
			continue
		}
	}
}
