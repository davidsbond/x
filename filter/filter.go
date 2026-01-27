// Package filter provides primitives for filtering slices of values.
package filter

type (
	// The Filter type is a function that returns true or false based on a given input.
	Filter[T any] func(T) bool
)

// All filters the slice of values down to only elements where each provided filter returned true. If no filters are
// provided, the slice is returned unchanged.
func All[T any](values []T, filters ...Filter[T]) []T {
	if len(filters) == 0 {
		return values
	}

	out := make([]T, 0, len(values))

	for _, value := range values {
		ok := true
		for _, filter := range filters {
			if !filter(value) {
				ok = false
				break
			}
		}

		if ok {
			out = append(out, value)
		}
	}

	return out
}

// Any filters the slice of values down to elements where at least one of the provided filters returned true.  If no
// filters are provided, the slice is returned unchanged.
func Any[T any](values []T, filters ...Filter[T]) []T {
	if len(filters) == 0 {
		return values
	}

	out := make([]T, 0, len(values))

	for _, value := range values {
		for _, filter := range filters {
			if filter(value) {
				out = append(out, value)
				break
			}
		}
	}

	return out
}
