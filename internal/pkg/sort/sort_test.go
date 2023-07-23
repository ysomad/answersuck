package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newSort(t *testing.T) {
	type args struct {
		col   string
		order string
	}
	tests := []struct {
		name string
		args args
		want Sort
	}{
		{
			name: "asc order",
			args: args{
				col:   "test",
				order: "asc",
			},
			want: Sort{
				col:   "test",
				order: "ASC",
			},
		},
		{
			name: "desc order",
			args: args{
				col:   "test",
				order: "desc",
			},
			want: Sort{
				col:   "test",
				order: "DESC",
			},
		},
		{
			name: "empty order",
			args: args{
				col:   "test",
				order: "",
			},
			want: Sort{
				col:   "test",
				order: "ASC",
			},
		},
		{
			name: "invalid order",
			args: args{
				col:   "test",
				order: "test",
			},
			want: Sort{
				col:   "test",
				order: "ASC",
			},
		},
		{
			name: "mixed asc order",
			args: args{
				col:   "test",
				order: "AsC",
			},
			want: Sort{
				col:   "test",
				order: "ASC",
			},
		},
		{
			name: "mixed desc order",
			args: args{
				col:   "test",
				order: "dEsC",
			},
			want: Sort{
				col:   "test",
				order: "DESC",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newSort(tt.args.col, tt.args.order)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewSortList(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    []Sort
		wantErr bool
	}{
		{
			name: "success mixed string",
			args: args{
				s: "field1,field2 ASC,field3 asc,field4 AsC,field5 DESC,field6 desc,field7 DeSc,field8",
			},
			want: []Sort{
				{
					col:   "field1",
					order: "ASC",
				},
				{
					col:   "field2",
					order: "ASC",
				},
				{
					col:   "field3",
					order: "ASC",
				},
				{
					col:   "field4",
					order: "ASC",
				},
				{
					col:   "field5",
					order: "DESC",
				},
				{
					col:   "field6",
					order: "DESC",
				},
				{
					col:   "field7",
					order: "DESC",
				},
				{
					col:   "field8",
					order: "ASC",
				},
			},
			wantErr: false,
		},
		{
			name: "success without orders",
			args: args{
				s: "field1,field2,field3,field4",
			},
			want: []Sort{
				{
					col:   "field1",
					order: "ASC",
				},
				{
					col:   "field2",
					order: "ASC",
				},
				{
					col:   "field3",
					order: "ASC",
				},
				{
					col:   "field4",
					order: "ASC",
				},
			},
			wantErr: false,
		},
		{
			name: "success 1 field",
			args: args{
				s: "field1",
			},
			want: []Sort{
				{
					col:   "field1",
					order: "ASC",
				},
			},
			wantErr: false,
		},
		{
			name: "empty sort string",
			args: args{
				s: "",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "invalid sort string",
			args: args{
				s: "field1 field2 field 3, field3 desc, field4 asc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "trailing spaces",
			args: args{
				s: " field1, field2, field3, field4 DESC, field5 AsC, field6 ",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSortList(tt.args.s)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}
