// Package channels provides functions for doing esoteric things with channels.
package channels

import (
	"context"
	"sync"
)

// Split writes all messages from the "from" channel to the subsequent "to" channels. Each "to" channel will receive
// every message sent on the "from" channel. This function blocks until the provided context is cancelled or the "from"
// channel is closed.
func Split[T any](ctx context.Context, from <-chan T, to ...chan<- T) {
	var group sync.WaitGroup

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-from:
			if !ok {
				return
			}

			for _, destination := range to {
				group.Add(1)
				go send(ctx, &group, destination, message)
			}

			group.Wait()
		}
	}
}

func send[T any](ctx context.Context, group *sync.WaitGroup, to chan<- T, message T) {
	defer group.Done()

	select {
	case <-ctx.Done():
		return
	case to <- message:
		return
	}
}

// Combine reads all messages from the "from" channels and writes them to the "to" channel. Ordering is not guaranteed
// at any point. This method blocks until the provided context is cancelled or all "from" channels are closed.
func Combine[T any](ctx context.Context, to chan<- T, from ...<-chan T) {
	var group sync.WaitGroup
	defer group.Wait()

	for _, source := range from {
		group.Add(1)
		go receive(ctx, &group, source, to)
	}
}

func receive[T any](ctx context.Context, group *sync.WaitGroup, from <-chan T, to chan<- T) {
	defer group.Done()

	for message := range from {
		select {
		case <-ctx.Done():
			return
		case to <- message:
			continue
		}
	}
}

type (
	// The Transformer type is a function that takes a single type as a parameter and returns another. To be used
	// with the Transform function.
	Transformer[A, B any] func(A) B
)

// Transform maps messages of a specified type read from the "from" channel. Messages are converted via the Transformer
// and written to the "to" channel. This function blocks until the provided context is cancelled or the "from" channel
// is closed.
func Transform[A, B any](ctx context.Context, from <-chan A, to chan<- B, fn Transformer[A, B]) {
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-from:
			if !ok {
				return
			}

			select {
			case <-ctx.Done():
				return
			case to <- fn(message):
				continue
			}
		}
	}
}

// Collect reads messages from the specified channel and collects them as a slice of the same type. This function
// blocks until the provided context is cancelled or the provided channel is closed.
func Collect[T any](ctx context.Context, from <-chan T) []T {
	values := make([]T, 0)

	for {
		select {
		case <-ctx.Done():
			return values
		case message, ok := <-from:
			if !ok {
				return values
			}

			values = append(values, message)
		}
	}
}
