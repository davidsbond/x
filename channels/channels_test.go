package channels_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/x/channels"
)

func TestSplit(t *testing.T) {
	t.Parallel()

	in := make(chan int, 1)

	out1 := make(chan int, 1)
	out2 := make(chan int, 1)
	out3 := make(chan int, 1)

	defer close(in)
	defer close(out1)
	defer close(out2)
	defer close(out3)

	go channels.Split(t.Context(), in, out1, out2, out3)

	in <- 42
	in <- 43
	in <- 44

	assert.Equal(t, 42, <-out1)
	assert.Equal(t, 42, <-out2)
	assert.Equal(t, 42, <-out3)

	assert.Equal(t, 43, <-out1)
	assert.Equal(t, 43, <-out2)
	assert.Equal(t, 43, <-out3)

	assert.Equal(t, 44, <-out1)
	assert.Equal(t, 44, <-out2)
	assert.Equal(t, 44, <-out3)
}

func TestCombine(t *testing.T) {
	t.Parallel()

	in1 := make(chan int, 1)
	in2 := make(chan int, 1)
	in3 := make(chan int, 1)

	out := make(chan int, 1)

	defer close(out)
	defer close(in1)
	defer close(in2)
	defer close(in3)

	go channels.Combine(t.Context(), out, in1, in2, in3)

	in1 <- 42
	in2 <- 43
	in3 <- 44

	results := make([]int, 3)
	for i := 0; i < 3; i++ {
		results[i] = <-out
	}

	assert.Contains(t, results, 42)
	assert.Contains(t, results, 43)
	assert.Contains(t, results, 44)
}

func TestTransform(t *testing.T) {
	t.Parallel()

	in := make(chan int, 1)
	out := make(chan string, 1)

	defer close(in)
	defer close(out)

	go channels.Transform(t.Context(), in, out, func(i int) string {
		return strconv.Itoa(i)
	})

	in <- 42
	assert.Equal(t, "42", <-out)
}

func TestCollect(t *testing.T) {
	t.Parallel()

	in := make(chan int, 3)
	for i := 0; i < 3; i++ {
		in <- i
	}

	close(in)
	out := channels.Collect(t.Context(), in)

	assert.Len(t, out, 3)
}
