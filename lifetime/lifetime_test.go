package lifetime_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidsbond/x/lifetime"
)

type (
	mockCloser struct {
		err    error
		closed bool
	}
)

func (mc *mockCloser) Close() error {
	mc.closed = true
	return mc.err
}

func TestLifetime(t *testing.T) {
	t.Parallel()

	t.Run("value is accessible within lifetime", func(t *testing.T) {
		expected := &mockCloser{}
		lt := lifetime.New(expected, time.Minute)

		actual, err := lt.Value()
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.False(t, expected.closed)
	})

	t.Run("value can be manually expired", func(t *testing.T) {
		expected := &mockCloser{}
		lt := lifetime.New(expected, time.Minute)
		lt.Expire()

		actual, err := lt.Value()
		require.Error(t, err)
		require.Nil(t, actual)
		assert.True(t, expected.closed)
	})

	t.Run("value expires after timeout", func(t *testing.T) {
		expected := &mockCloser{}
		lt := lifetime.New(expected, time.Second/2)

		<-time.After(time.Second)
		actual, err := lt.Value()
		require.Error(t, err)
		require.Nil(t, actual)
		assert.True(t, expected.closed)
	})

	t.Run("lifetime returns any close errors", func(t *testing.T) {
		expected := &mockCloser{err: io.EOF}
		lt := lifetime.New(expected, time.Minute)
		lt.Expire()

		actual, err := lt.Value()
		require.Error(t, err)
		require.Nil(t, actual)
		assert.True(t, expected.closed)
		assert.True(t, errors.Is(err, io.EOF))
	})

	t.Run("lifetime can be reset", func(t *testing.T) {
		expected := &mockCloser{}
		lt := lifetime.New(expected, time.Minute)

		actual, err := lt.Value()
		require.NoError(t, err)
		require.NotNil(t, actual)
		assert.False(t, expected.closed)

		require.NoError(t, lt.Reset(time.Second/2))
		<-time.After(time.Second)

		actual, err = lt.Value()
		require.Error(t, err)
		require.Nil(t, actual)
		assert.True(t, expected.closed)
	})

	t.Run("can't call reset on expired lifetime", func(t *testing.T) {
		expected := &mockCloser{}
		lt := lifetime.New(expected, time.Second/2)

		<-time.After(time.Second)
		require.Error(t, lt.Reset(time.Second))
		assert.True(t, expected.closed)
	})
}
