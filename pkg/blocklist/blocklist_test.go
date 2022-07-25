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

	for _, tt := range []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "found",
			value: "admin",
			want:  true,
		},
		{
			name:  "not found",
			value: "ysomad",
			want:  false,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			found := bl.Find(tt.value)
			assert.Equal(t, found, tt.want)
		})
	}
}
