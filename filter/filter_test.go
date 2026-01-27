package filter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidsbond/x/filter"
)

func TestAll(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name     string
		Input    []string
		Filters  []filter.Filter[string]
		Expected []string
	}{
		{
			Name:     "all elements when no filters",
			Input:    []string{"a", "b", "c"},
			Expected: []string{"a", "b", "c"},
		},
		{
			Name:     "filters elements that match",
			Input:    []string{"a", "b", "c"},
			Expected: []string{"a", "c"},
			Filters: []filter.Filter[string]{
				func(s string) bool {
					return s != "b"
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			actual := filter.All(tc.Input, tc.Filters...)
			assert.EqualValues(t, tc.Expected, actual)
		})
	}
}

func TestAny(t *testing.T) {
	t.Parallel()

	tt := []struct {
		Name     string
		Input    []string
		Filters  []filter.Filter[string]
		Expected []string
	}{
		{
			Name:     "all elements when no filters",
			Input:    []string{"a", "b", "c"},
			Expected: []string{"a", "b", "c"},
		},
		{
			Name:     "filters elements that match",
			Input:    []string{"a", "b", "c"},
			Expected: []string{"b", "c"},
			Filters: []filter.Filter[string]{
				func(s string) bool {
					return s == "b"
				},
				func(s string) bool {
					return s == "c"
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			actual := filter.Any(tc.Input, tc.Filters...)
			assert.EqualValues(t, tc.Expected, actual)
		})
	}
}
