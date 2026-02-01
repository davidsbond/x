// Package convert provides functions for converting types from one to another.
package convert

// Slice maps a slice of one type into a slice of another using the provided conversion function.
func Slice[Have, Want any](in []Have, fn func(Have) Want) []Want {
	out := make([]Want, len(in))
	for i := range in {
		out[i] = fn(in[i])
	}

	return out
}

// Map consumes a map, running each key and value through a conversion function, returning a new map with the same
// keys but transformed values.
func Map[Key comparable, Have any, Want any](in map[Key]Have, fn func(Key, Have) Want) map[Key]Want {
	out := make(map[Key]Want, len(in))
	for k, v := range in {
		out[k] = fn(k, v)
	}

	return out
}
