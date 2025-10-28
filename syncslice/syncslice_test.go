package syncslice_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/x/syncslice"
)

func Test_New(t *testing.T) {
	t.Parallel()

	slice := syncslice.New[string]()

	t.Run("appends", func(t *testing.T) {
		slice.Append("a", "b", "c")
		require.Equal(t, 3, slice.Len())
	})

	t.Run("indexes", func(t *testing.T) {
		a := slice.At(0)
		b := slice.At(1)
		c := slice.At(2)

		require.Equal(t, "a", a)
		require.Equal(t, "b", b)
		require.Equal(t, "c", c)

		slice.Set(0, "a1")
		require.Equal(t, "a1", slice.At(0))
	})

	t.Run("ranges", func(t *testing.T) {
		for i, exp := range slice.Range() {
			act := slice.At(i)

			require.Equal(t, exp, act)
		}
	})

	t.Run("unwraps", func(t *testing.T) {
		act := slice.Unwrap()

		require.NotNil(t, act)
		require.Equal(t, len(act), slice.Len())
	})
}

func Test_NewLen(t *testing.T) {
	t.Parallel()

	slice := syncslice.NewLen[string](10)

	assert.Equal(t, 10, slice.Len())
}

func Test_NewLenCap(t *testing.T) {
	t.Parallel()

	slice := syncslice.NewLenCap[string](5, 10)

	assert.Equal(t, 5, slice.Len())
	assert.Equal(t, 10, slice.Cap())
}
