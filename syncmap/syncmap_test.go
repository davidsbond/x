package syncmap_test

import (
	"strconv"
	"testing"

	"github.com/davidsbond/x/syncmap"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	t.Parallel()

	m := syncmap.New[string, int]()

	t.Run("adds entries", func(t *testing.T) {
		const (
			k = "a"
			v = 1
		)

		m.Put(k, v)
		require.Equal(t, 1, m.Len())
		actual, ok := m.Get(k)
		require.True(t, ok)
		require.Equal(t, v, actual)
	})

	t.Run("lists keys", func(t *testing.T) {
		keys := m.Keys()

		require.Len(t, keys, m.Len())
		require.Contains(t, keys, "a")
	})

	t.Run("lists values", func(t *testing.T) {
		values := m.Values()

		require.Len(t, values, m.Len())
		require.Contains(t, values, 1)
	})

	t.Run("ranges", func(t *testing.T) {
		count := 0
		for k, v := range m.Range() {
			require.Equal(t, "a", k)
			require.Equal(t, v, 1)
			count++
		}

		require.Equal(t, m.Len(), count)
	})

	t.Run("removes entries", func(t *testing.T) {
		m.Remove("a")
		require.Zero(t, m.Len())

		v, ok := m.Get("a")
		require.False(t, ok)
		require.Zero(t, v)
	})

	t.Run("clears all entries", func(t *testing.T) {
		for i := range 10 {
			m.Put(strconv.Itoa(i), i)
		}

		m.Clear()
		require.Zero(t, m.Len())
	})
}
