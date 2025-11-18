package envvar_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/x/envvar"
)

func TestInt(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		t.Setenv("TEST_INT", "42")
		actual := envvar.Int("TEST_INT", 0)
		assert.EqualValues(t, 42, actual)
	})

	t.Run("returns default", func(t *testing.T) {
		t.Setenv("TEST_INT", "")
		actual := envvar.Int("TEST_INT", 32)
		assert.EqualValues(t, 32, actual)
	})

	t.Run("handles invalid", func(t *testing.T) {
		t.Setenv("TEST_INT", "abc")
		actual := envvar.Int("TEST_INT", 32)
		assert.EqualValues(t, 32, actual)
	})
}

func TestInt64(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		t.Setenv("TEST_INT", "42")
		actual := envvar.Int64("TEST_INT", 0)
		assert.EqualValues(t, 42, actual)
	})

	t.Run("returns default", func(t *testing.T) {
		t.Setenv("TEST_INT", "")
		actual := envvar.Int64("TEST_INT", 32)
		assert.EqualValues(t, 32, actual)
	})

	t.Run("handles invalid", func(t *testing.T) {
		t.Setenv("TEST_INT", "abc")
		actual := envvar.Int64("TEST_INT", 32)
		assert.EqualValues(t, 32, actual)
	})
}

func TestFloat64(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		t.Setenv("TEST_FLOAT", "42.1")
		actual := envvar.Float64("TEST_FLOAT", 0)
		assert.EqualValues(t, 42.1, actual)
	})

	t.Run("returns default", func(t *testing.T) {
		t.Setenv("TEST_FLOAT", "")
		actual := envvar.Float64("TEST_FLOAT", 32.2)
		assert.EqualValues(t, 32.2, actual)
	})

	t.Run("handles invalid", func(t *testing.T) {
		t.Setenv("TEST_FLOAT", "abc")
		actual := envvar.Float64("TEST_FLOAT", 32.2)
		assert.EqualValues(t, 32.2, actual)
	})
}

func TestUint64(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		t.Setenv("TEST_INT", "42")
		actual := envvar.Uint64("TEST_INT", 0)
		assert.EqualValues(t, 42, actual)
	})

	t.Run("returns default", func(t *testing.T) {
		t.Setenv("TEST_INT", "")
		actual := envvar.Uint64("TEST_INT", 32)
		assert.EqualValues(t, 32, actual)
	})

	t.Run("handles invalid", func(t *testing.T) {
		t.Setenv("TEST_INT", "abc")
		actual := envvar.Uint64("TEST_INT", 32)
		assert.EqualValues(t, 32, actual)
	})
}

func TestBool(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		t.Setenv("TEST_BOOL", "true")
		actual := envvar.Bool("TEST_BOOL", false)
		assert.True(t, actual)

		t.Setenv("TEST_BOOL", "false")
		actual = envvar.Bool("TEST_BOOL", true)
		assert.False(t, actual)
	})

	t.Run("returns default", func(t *testing.T) {
		t.Setenv("TEST_BOOL", "")
		actual := envvar.Bool("TEST_BOOL", true)
		assert.True(t, actual)
	})

	t.Run("handles invalid", func(t *testing.T) {
		t.Setenv("TEST_BOOL", "abc")
		actual := envvar.Bool("TEST_BOOL", true)
		assert.True(t, actual)
	})
}

func TestString(t *testing.T) {
	t.Run("parses correctly", func(t *testing.T) {
		t.Setenv("TEST_STRING", "foo")
		actual := envvar.String("TEST_STRING", "")
		assert.EqualValues(t, "foo", actual)
	})

	t.Run("returns default", func(t *testing.T) {
		t.Setenv("TEST_STRING", "")
		actual := envvar.String("TEST_STRING", "foo")
		assert.EqualValues(t, "foo", actual)
	})
}
