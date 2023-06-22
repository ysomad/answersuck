package paging

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

func TestNewIntSeek(t *testing.T) {
	type args[T constraints.Signed] struct {
		lastID   T
		pageSize int32
	}

	type test[T constraints.Signed] struct {
		name string
		args args[T]
		want IntSeek[T]
	}

	tests := []test[int32]{
		{
			name: "success",
			args: args[int32]{
				lastID:   69,
				pageSize: 50,
			},
			want: IntSeek[int32]{
				PageSize: 50,
				LastID:   69,
			},
		},
		{
			name: "empty lastID",
			args: args[int32]{
				pageSize: 50,
			},
			want: IntSeek[int32]{
				PageSize: 50,
				LastID:   0,
			},
		},
		{
			name: "empty page size",
			args: args[int32]{
				lastID: 69,
			},
			want: IntSeek[int32]{
				LastID:   69,
				PageSize: defaultPageSize,
			},
		},
		{
			name: "page size less than min page size",
			args: args[int32]{
				lastID:   69,
				pageSize: 0,
			},
			want: IntSeek[int32]{
				LastID:   69,
				PageSize: defaultPageSize,
			},
		},
		{
			name: "page size greater than max page size",
			args: args[int32]{
				lastID:   69,
				pageSize: 1111,
			},
			want: IntSeek[int32]{
				LastID:   69,
				PageSize: maxPageSize,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIntSeek(tt.args.lastID, tt.args.pageSize)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewInfList(t *testing.T) {
	type args[T any] struct {
		items    []T
		pageSize int32
	}

	type test[T any] struct {
		name    string
		args    args[T]
		want    List[T]
		wantErr bool
	}

	tests := []test[string]{
		{
			name: "has next page 1",
			args: args[string]{
				items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
					"item 6",
				},
				pageSize: 5,
			},
			want: List[string]{
				Items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
				},
				HasNext: true,
			},
			wantErr: false,
		},
		{
			name: "has next page 2",
			args: args[string]{
				items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
					"item 6",
					"item 7",
					"item 8",
					"item 9",
					"item 10",
					"item 11",
				},
				pageSize: 10,
			},
			want: List[string]{
				Items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
					"item 6",
					"item 7",
					"item 8",
					"item 9",
					"item 10",
				},
				HasNext: true,
			},
			wantErr: false,
		},
		{
			name: "does not have next page 1",
			args: args[string]{
				items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
				},
				pageSize: 5,
			},
			want: List[string]{
				Items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
				},
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name: "does not have next page 2",
			args: args[string]{
				items: []string{
					"item 1",
					"item 2",
					"item 3",
				},
				pageSize: 5,
			},
			want: List[string]{
				Items: []string{
					"item 1",
					"item 2",
					"item 3",
				},
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name: "does not have next page 2",
			args: args[string]{
				items: []string{
					"item 1",
					"item 2",
					"item 3",
				},
				pageSize: 5,
			},
			want: List[string]{
				Items: []string{
					"item 1",
					"item 2",
					"item 3",
				},
				HasNext: false,
			},
			wantErr: false,
		},
		{
			name: "items has length more than pageSize + 1",
			args: args[string]{
				items: []string{
					"item 1",
					"item 2",
					"item 3",
					"item 4",
					"item 5",
					"item 6",
					"item 7",
					"item 8",
					"item 9",
					"item 10",
					"item 11",
					"item 12",
				},
				pageSize: 10,
			},
			want:    List[string]{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewList(tt.args.items, tt.args.pageSize)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, (err != nil))
		})
	}
}

var _ ListItem = &testList{}

type testList struct {
	id        string
	createdAt time.Time
}

func (l *testList) GetID() string      { return l.id }
func (l *testList) GetTime() time.Time { return l.createdAt }

func TestNewTokenList(t *testing.T) {
	type args[T ListItem] struct {
		items    []T
		pageSize int32
	}
	type test[T ListItem] struct {
		name    string
		args    args[T]
		want    TokenList[T]
		wantErr bool
	}
	time1 := time.Time{}.Add(time.Hour)
	time2 := time1.Add(time.Hour)
	time3 := time2.Add(time.Hour)
	time4 := time3.Add(time.Hour)
	time5 := time4.Add(time.Hour)
	lastItemTestTime := time.Time{}.Add(time.Hour * 99999)
	tests := []test[*testList]{
		{
			name: "0 items",
			args: args[*testList]{
				items:    []*testList{},
				pageSize: 5,
			},
			want: TokenList[*testList]{
				Items: []*testList{},
			},
			wantErr: false,
		},
		{
			name: "6 items, has next page",
			args: args[*testList]{
				items: []*testList{
					{
						id:        "id1",
						createdAt: time1,
					},
					{
						id:        "id2",
						createdAt: time2,
					},
					{
						id:        "id3",
						createdAt: time3,
					},
					{
						id:        "id4",
						createdAt: time4,
					},
					{
						id:        "id5",
						createdAt: time5,
					},
					{
						id:        "id6",
						createdAt: lastItemTestTime,
					},
				},
				pageSize: 5,
			},
			want: TokenList[*testList]{
				Items: []*testList{
					{
						id:        "id1",
						createdAt: time1,
					},
					{
						id:        "id2",
						createdAt: time2,
					},
					{
						id:        "id3",
						createdAt: time3,
					},
					{
						id:        "id4",
						createdAt: time4,
					},
					{
						id:        "id5",
						createdAt: time5,
					},
				},
				NextPageToken: "aWQ1LDAwMDEtMDEtMDFUMDU6MDA6MDBa",
			},
			wantErr: false,
		},
		{
			name: "5 items, no next page",
			args: args[*testList]{
				items: []*testList{
					{
						id:        "id1",
						createdAt: time1,
					},
					{
						id:        "id2",
						createdAt: time2,
					},
					{
						id:        "id3",
						createdAt: time3,
					},
					{
						id:        "id4",
						createdAt: time4,
					},
					{
						id:        "id5",
						createdAt: lastItemTestTime,
					},
				},
				pageSize: 5,
			},
			want: TokenList[*testList]{
				Items: []*testList{
					{
						id:        "id1",
						createdAt: time1,
					},
					{
						id:        "id2",
						createdAt: time2,
					},
					{
						id:        "id3",
						createdAt: time3,
					},
					{
						id:        "id4",
						createdAt: time4,
					},
					{
						id:        "id5",
						createdAt: lastItemTestTime,
					},
				},
				NextPageToken: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTokenList(tt.args.items, tt.args.pageSize)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}
