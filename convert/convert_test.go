package convert_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/x/convert"
)

func TestSlice(t *testing.T) {
	t.Parallel()

	input := []int{1, 2, 3}
	expected := []string{"1", "2", "3"}

	actual := convert.Slice(input, func(value int) string {
		assert.Contains(t, input, value)

		return strconv.Itoa(value)
	})

	assert.EqualValues(t, expected, actual)
}

func TestMap(t *testing.T) {
	t.Parallel()

	input := map[string]int{"a": 1, "b": 2, "c": 3}
	expected := map[string]string{"a": "1", "b": "2", "c": "3"}

	actual := convert.Map(input, func(key string, value int) string {
		assert.Contains(t, input, key)

		return strconv.Itoa(value)
	})

	assert.EqualValues(t, expected, actual)
}
