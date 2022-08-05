package repository_psql

import (
	"context"
	"testing"

	"github.com/answersuck/host/internal/adapter/repository/psql"
	"github.com/answersuck/host/internal/domain/tag"
	"github.com/stretchr/testify/assert"
)

var _tagRepo *psql.TagRepo

func TestTagRepo_SaveMultiple(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		tags []tag.Tag
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "tags saved successfully",
			args: args{
				ctx: context.Background(),
				tags: []tag.Tag{
					{
						Name:       "tag1",
						LanguageId: 1,
					},
					{
						Name:       "tag2",
						LanguageId: 2,
					},
					{
						Name:       "tag3",
						LanguageId: 1,
					},
					{
						Name:       "tag3",
						LanguageId: 1,
					},
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "language not found",
			args: args{
				ctx: context.Background(),
				tags: []tag.Tag{
					{
						Name:       "tag1",
						LanguageId: 999999, // doesnt exist
					},
					{
						Name:       "tag2",
						LanguageId: 2,
					},
				},
			},
			wantErr: true,
			err:     tag.ErrLanguageIdNotFound,
		},
		{
			name: "empty list of tags",
			args: args{
				ctx:  context.Background(),
				tags: []tag.Tag{},
			},
			wantErr: true,
			err:     tag.ErrEmptyTagList,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _tagRepo.SaveMultiple(tt.args.ctx, tt.args.tags)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}

			assert.Equal(t, tt.wantErr, (err != nil))
			for _, tag := range got {
				assert.NotEmpty(t, tag.Id)
				assert.NotEmpty(t, tag.Name)
				assert.NotEmpty(t, tag.LanguageId)
			}
		})
	}
}

// TODO: finish when sorting, filtering and pagination is done
// func TestTagRepo_FindAll(t *testing.T) {
// 	t.Parallel()
//
// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []tag.Tag
// 		wantErr bool
// 		err     error
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// got, err := _tagRepo.FindAll(tt.args.ctx)
// 		})
// 	}
// }
