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
		slice := weightslice.New[string, int](weightslice.Descending)

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
		slice := weightslice.New[string, int](weightslice.Ascending)

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
}
