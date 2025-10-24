package slicepool_test

import (
	"testing"

	"github.com/davidsbond/x/slicepool"
	"github.com/stretchr/testify/require"
)

func TestPool(t *testing.T) {
	t.Parallel()

	pool := slicepool.New[byte](64)

	t.Run("provides slices", func(t *testing.T) {
		v := pool.Get()

		require.NotNil(t, v)
		require.Len(t, *v, 64)
	})

	t.Run("reuses slices", func(t *testing.T) {
		v := pool.Get()

		buf := *v
		buf[0] = 1

		pool.Put(v)

		v = pool.Get()

		buf = *v
		require.EqualValues(t, 1, buf[0])
	})
}
