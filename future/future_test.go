package future_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/x/future"
)

func TestDo(t *testing.T) {
	t.Parallel()

	t.Run("returns results", func(t *testing.T) {
		const expected = 42

		f := future.Do(t.Context(), func(ctx context.Context) (int, error) {
			return expected, nil
		})

		actual, err := f.Result(t.Context())
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})

	t.Run("returns error", func(t *testing.T) {
		f := future.Do(t.Context(), func(ctx context.Context) (int, error) {
			return 0, io.EOF
		})

		_, err := f.Result(t.Context())
		require.Equal(t, io.EOF, err)
	})

	t.Run("handles cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		f := future.Do(ctx, func(ctx context.Context) (int, error) {
			return 42, ctx.Err()
		})

		v, err := f.Result(ctx)
		require.Zero(t, v)
		require.Equal(t, context.Canceled, err)
	})
}

func TestAll(t *testing.T) {
	t.Parallel()

	t.Run("returns all results", func(t *testing.T) {
		results := future.All(t.Context(),
			func(ctx context.Context) (int, error) {
				return 42, nil
			},
			func(ctx context.Context) (int, error) {
				return 0, io.EOF
			},
			func(ctx context.Context) (int, error) {
				return 1, nil
			},
		)

		var i int
		for result, err := range results {
			if i == 0 {
				assert.Equal(t, 42, result)
				assert.NoError(t, err)
				i++
				continue
			}

			if i == 1 {
				assert.Zero(t, result)
				assert.Error(t, err)
				i++
				continue
			}

			if i == 2 {
				assert.Equal(t, 1, result)
				assert.NoError(t, err)
			}
		}

		assert.Equal(t, 2, i)
	})
}
