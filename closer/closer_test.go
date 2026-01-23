package closer_test

import (
	"database/sql"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/x/closer"
)

type (
	testCloser struct {
		closed bool
		err    error
	}
)

func (tc *testCloser) Close() error {
	tc.closed = true
	return tc.err
}

func TestCollection_Close(t *testing.T) {
	t.Parallel()

	t.Run("closes all closers", func(t *testing.T) {
		closers := []*testCloser{
			{},
			{},
			{},
		}

		collection := closer.NewCollection()
		for _, cl := range closers {
			collection.Add(cl)
		}

		err := collection.Close()
		require.NoError(t, err)

		for i, cl := range closers {
			assert.True(t, cl.closed, i)
		}
	})

	t.Run("returns all errors", func(t *testing.T) {
		closers := []*testCloser{
			{err: io.EOF},
			{err: net.ErrClosed},
			{err: sql.ErrNoRows},
		}

		collection := closer.NewCollection()
		for _, cl := range closers {
			collection.Add(cl)
		}

		err := collection.Close()
		require.Error(t, err)

		for i, cl := range closers {
			assert.True(t, cl.closed, i)
			assert.ErrorIs(t, err, cl.err, i)
		}
	})
}
