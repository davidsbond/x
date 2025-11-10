// Package future provides types and methods for performing asynchronous actions that have return values that can
// be waited upon.
package future

import (
	"context"
	"iter"
)

type (
	// The Future type is used to perform an asynchronous action and obtain its result at a later time.
	Future[T any] struct {
		signal chan result[T]
	}

	// The Func type represents the action the future should take upon creation.
	Func[T any] func(ctx context.Context) (T, error)

	result[T any] struct {
		value T
		err   error
	}
)

// Do returns a new Future that executes the given Func implementation. The function is invoked immediately within
// its own goroutine and its return value is made accessible via the Result method. The provided context is passed
// into the Func call for cancellation.
func Do[T any](ctx context.Context, fn Func[T]) *Future[T] {
	f := &Future[T]{
		signal: make(chan result[T], 1),
	}

	go func() {
		defer close(f.signal)
		value, err := fn(ctx)
		f.signal <- result[T]{value: value, err: err}
	}()

	return f
}

// All creates a Future for each provided Func, running each in its own goroutine. Results of each Future are exposed
// via an iter.Seq2 that can be ranged over. The first value of the iter.Seq2 is the return type of the Future, the
// second is any error that has occurred.
func All[T any](ctx context.Context, funcs ...Func[T]) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		futures := make([]*Future[T], len(funcs))
		for i, f := range funcs {
			futures[i] = Do(ctx, f)
		}

		for _, future := range futures {
			if !yield(future.Result(ctx)) {
				return
			}
		}
	}
}

// Result returns the result of the Future. If the invoked Func is still in-progress, this method blocks until the
// Func has returned or the provided context is cancelled.
func (f *Future[T]) Result(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		var v T
		return v, ctx.Err()
	case r := <-f.signal:
		return r.value, r.err
	}
}
