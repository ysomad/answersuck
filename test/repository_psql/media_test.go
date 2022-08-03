package repository_psql

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/answersuck/host/internal/adapter/repository/psql"
	"github.com/answersuck/host/internal/domain/account"
	"github.com/answersuck/host/internal/domain/media"
)

var _mediaRepo *psql.MediaRepo

func TestMediaRepo_Save(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	type args struct {
		ctx context.Context
		m   media.Media
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "media saved",
			args: args{
				ctx: context.Background(),
				m: media.Media{
					Id:        "76db5e01-deee-414f-b4ee-8e649fd372b2",
					Filename:  "filename",
					Type:      media.TypeAudioMP4,
					AccountId: a.Id,
					CreatedAt: time.Now(),
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "media already exist",
			args: args{
				ctx: context.Background(),
				m: media.Media{
					Id:        "76db5e01-deee-414f-b4ee-8e649fd372b2",
					Type:      media.TypeAudioMP4,
					AccountId: a.Id,
				},
			},
			wantErr: true,
			err:     media.ErrAlreadyExist,
		},
		{
			name: "account not found",
			args: args{
				ctx: context.Background(),
				m: media.Media{
					Id:        "76db5e01-deee-414f-b4ee-8e649fd372b3",
					Type:      media.TypeAudioMP4,
					AccountId: "76db5e01-deee-414f-b4ee-8e649fd372b2",
				},
			},
			wantErr: true,
			err:     media.ErrAccountNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := _mediaRepo.Save(tt.args.ctx, tt.args.m)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			var got media.Media
			err = _mediaRepo.Pool.QueryRow(
				context.Background(),
				"SELECT filename, type, account_id, created_at FROM media WHERE id = $1",
				tt.args.m.Id,
			).Scan(&got.Filename, &got.Type, &got.AccountId, &got.CreatedAt)
			assert.NoError(t, err)

			assert.Equal(t, tt.args.m.Filename, got.Filename)
			assert.Equal(t, tt.args.m.Type, got.Type)
			assert.Equal(t, tt.args.m.AccountId, got.AccountId)
			assert.Equal(t, tt.args.m.CreatedAt.Unix(), got.CreatedAt.Unix())
		})
	}
}

//
// func Test_mediaRepo_FindMediaTypeById(t *testing.T) {
// 	type fields struct {
// 		Logger *zap.Logger
// 		Client *postgres.Client
// 	}
// 	type args struct {
// 		ctx     context.Context
// 		mediaId string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    string
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &mediaRepo{
// 				Logger: tt.fields.Logger,
// 				Client: tt.fields.Client,
// 			}
// 			got, err := r.FindMediaTypeById(tt.args.ctx, tt.args.mediaId)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("mediaRepo.FindMediaTypeById() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("mediaRepo.FindMediaTypeById() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
