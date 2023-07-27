package paging

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewListWithOffset(t *testing.T) {
	type args[T any] struct {
		items  []T
		limit  uint64
		offset uint64
	}
	tests := []struct {
		name    string
		args    args[int]
		want    List[int]
		wantErr bool
	}{
		{
			name: "has no next page",
			args: args[int]{
				items:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				limit:  10,
				offset: 0,
			},
			want: List[int]{
				Items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				NextPageToken: "",
			},
			wantErr: false,
		},
		{
			name: "has next page",
			args: args[int]{
				items:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
				limit:  10,
				offset: 0,
			},
			want: List[int]{
				Items:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				NextPageToken: "CgAAAAAAAAAKAAAAAAAAAA==", // limit = 10, offset = 10
			},
			wantErr: false,
		},
		{
			name: "invalid args 1",
			args: args[int]{
				items:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
				limit:  10,
				offset: 0,
			},
			want:    List[int]{},
			wantErr: true,
		},
		{
			name: "empty items",
			args: args[int]{
				items:  []int(nil),
				limit:  10,
				offset: 0,
			},
			want:    List[int]{},
			wantErr: false,
		},
		{
			name: "short list",
			args: args[int]{
				items:  []int{1, 2, 3},
				limit:  2,
				offset: 0,
			},
			want: List[int]{
				Items:         []int{1, 2},
				NextPageToken: "AgAAAAAAAAACAAAAAAAAAA==",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewListWithOffset(tt.args.items, tt.args.limit, tt.args.offset)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}
