// Package convert provides functions for converting types from one to another.
package convert

// Slice maps a slice of one type into a slice of another using the provided conversion function.
func Slice[In, Out any](in []In, fn func(In) Out) []Out {
	out := make([]Out, len(in))
	for i := range in {
		out[i] = fn(in[i])
	}

	return out
}
