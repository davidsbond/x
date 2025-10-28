package syncset_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/davidsbond/x/syncset"
)

func TestSet(t *testing.T) {
	t.Parallel()

	m := syncset.New[string]()

	t.Run("adds entries", func(t *testing.T) {
		const (
			k = "a"
		)

		m.Put(k)
		require.Equal(t, 1, m.Len())
	})

	t.Run("lists values", func(t *testing.T) {
		values := m.Values()

		require.Len(t, values, m.Len())
		require.Contains(t, values, "a")
	})

	t.Run("ranges", func(t *testing.T) {
		count := 0
		for v := range m.Range() {
			require.Equal(t, "a", v)
			count++
		}

		require.Equal(t, m.Len(), count)
	})

	t.Run("removes entries", func(t *testing.T) {
		m.Remove("a")
		require.Zero(t, m.Len())
	})

	t.Run("clears all entries", func(t *testing.T) {
		for i := range 10 {
			m.Put(strconv.Itoa(i))
		}

		m.Clear()
		require.Zero(t, m.Len())
	})
}
