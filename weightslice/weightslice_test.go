package weightslice_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/x/weightslice"
)

func TestWeightSlice(t *testing.T) {
	t.Parallel()

	t.Run("descending", func(t *testing.T) {
		slice := weightslice.New[string, int](nil, weightslice.Descending)

		for i := 0; i < 10; i++ {
			// Add elements with weights equal to their index.
			slice.Append(strconv.Itoa(i), i)
		}

		for i, v := range slice.Range() {
			// First element should be the highest.
			assert.Equal(t, "9", v)

			// Modify the weight to make it the lowest value, next we'll check if it's the last element.
			slice.SetWeight(i, -100)
			break
		}

		for _, v := range slice.Range() {
			// 8 Should now be the first element.
			assert.Equal(t, "8", v)
			break
		}
	})

	t.Run("ascending", func(t *testing.T) {
		slice := weightslice.New[string, int](nil, weightslice.Ascending)

		for i := 0; i < 10; i++ {
			// Add elements with weights equal to their index.
			slice.Append(strconv.Itoa(i), i)
		}

		for i, v := range slice.Range() {
			// First element should be the highest.
			assert.Equal(t, "0", v)

			// Modify the weight to make it the lowest value, next we'll check if it's the last element.
			slice.SetWeight(i, 100)
			break
		}

		for _, v := range slice.Range() {
			// 1 Should now be the first element.
			assert.Equal(t, "1", v)
			break
		}
	})

	t.Run("initial slice", func(t *testing.T) {
		slice := weightslice.New[string, int]([]string{"a", "b", "c"}, weightslice.Ascending)

		for i, v := range slice.Range() {
			// "a" should be the first element
			assert.Equal(t, "a", v)

			// Modify the weight of "a" so it comes last
			slice.SetWeight(i, 100)
			break
		}

		for _, v := range slice.Range() {
			// "b" should now be the first element
			assert.Equal(t, "b", v)
			break
		}
	})

	t.Run("resets slice", func(t *testing.T) {
		slice := weightslice.New[string, int]([]string{"a", "b", "c"}, weightslice.Ascending)

		for i, v := range slice.Range() {
			// "a" should be the first element
			assert.Equal(t, "a", v)

			// Modify the weight of "a" so it comes last
			slice.SetWeight(i, 100)
			break
		}

		// We don't have a mechanism for introspecting the weights for now so this test pretty
		// much just makes sure we don't block indefinitely or panic.
		slice.Reset()
	})
}
