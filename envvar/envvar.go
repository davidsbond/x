// Package envvar provides functions for working with environment variables. Primarily with parsing their string
// values into different types and managing default values.
package envvar

import (
	"os"
	"strconv"
)

// String returns the value of the specified environment variable as a string. Returns the specified default value when
// the environment variable is not set.
func String(key string, def string) string {
	return parse(key, def, func(s string) (string, error) {
		return s, nil
	})
}

// Int returns the value of the specified environment variable as an integer. Returns the specified default value when
// the environment variable is not set or cannot be parsed.
func Int(key string, def int) int {
	return parse(key, def, strconv.Atoi)
}

// Int64 returns the value of the specified environment variable as a 64-bit integer. Returns the specified default
// value when the environment variable is not set or cannot be parsed.
func Int64(key string, def int64) int64 {
	return parse(key, def, func(s string) (int64, error) {
		return strconv.ParseInt(s, 10, 64)
	})
}

// Bool returns the value of the specified environment variable as a boolean. Accepts "true" and "false". Returns the
// specified default value when the environment variable is not set or cannot be parsed.
func Bool(key string, def bool) bool {
	return parse(key, def, strconv.ParseBool)
}

// Float64 returns the value of the specified environment variable as a 64-bit floating point number. Returns the
// specified default value when the environment variable is not set or cannot be parsed.
func Float64(key string, def float64) float64 {
	return parse(key, def, func(s string) (float64, error) {
		return strconv.ParseFloat(s, 64)
	})
}

// Uint64 returns the value of the specified environment variable as an unsigned 64-bit integer. Returns the specified
// default value when the environment variable is not set or cannot be parsed.
func Uint64(key string, def uint64) uint64 {
	return parse(key, def, func(s string) (uint64, error) {
		return strconv.ParseUint(s, 10, 64)
	})
}

func parse[T any](key string, def T, parser func(string) (T, error)) T {
	str := os.Getenv(key)
	if str == "" {
		return def
	}

	value, err := parser(str)
	if err != nil {
		return def
	}

	return value
}
