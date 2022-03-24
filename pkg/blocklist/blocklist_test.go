package blocklist

import (
	"github.com/go-playground/assert/v2"
	"sort"
	"testing"
)

func TestNew(t *testing.T) {
	bl := New(WithUsernames)

	assert.Equal(t, sort.StringsAreSorted(bl.values), true)
	assert.Equal(t, len(WithUsernames()), len(bl.values))
}

func TestFind(t *testing.T) {
	type testCase struct {
		description string
		value       string
		expected    bool
	}

	bl := New(WithUsernames)

	for _, c := range []testCase{
		{
			description: "found",
			value:       "admin",
			expected:    true,
		},
		{
			description: "not found",
			value:       "ysomad",
			expected:    false,
		},
	} {
		t.Run(c.description, func(t *testing.T) {
			found := bl.Find(c.value)
			assert.Equal(t, found, c.expected)
		})
	}
}
