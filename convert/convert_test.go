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

	actual := convert.Slice(input, func(in int) string {
		return strconv.Itoa(in)
	})

	assert.EqualValues(t, expected, actual)
}
