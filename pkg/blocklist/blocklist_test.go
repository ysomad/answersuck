package blocklist

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	bl := New(WithUsernames)

	assert.Equal(t, sort.StringsAreSorted(bl.values), true)
	assert.Equal(t, len(WithUsernames()), len(bl.values))
}

func Test_Find(t *testing.T) {
	bl := New(WithUsernames)

	for _, c := range []struct {
		description string
		value       string
		want        bool
	}{
		{
			description: "found",
			value:       "admin",
			want:        true,
		},
		{
			description: "not found",
			value:       "ysomad",
			want:        false,
		},
	} {
		t.Run(c.description, func(t *testing.T) {
			found := bl.Find(c.value)
			assert.Equal(t, found, c.want)
		})
	}
}
