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

func insertTestMedia(m *media.Media) (*media.Media, error) {
	m.CreatedAt = time.Now()
	m.Filename = "test"
	if m.Type == "" {
		m.Type = media.TypeImageJPEG
	}
	if err := _mediaRepo.Pool.QueryRow(
		context.Background(),
		"INSERT INTO media(filename, type, account_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		m.Filename, m.Type, m.AccountId, m.CreatedAt,
	).Scan(&m.Id); err != nil {
		return nil, err
	}
	return m, nil
}

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

func Test_mediaRepo_FindMediaTypeById(t *testing.T) {
	t.Parallel()

	a, err := insertTestAccount(account.Account{})
	assert.NoError(t, err)

	m, err := insertTestMedia(&media.Media{AccountId: a.Id})
	assert.NoError(t, err)

	type args struct {
		ctx     context.Context
		mediaId string
	}
	tests := []struct {
		name    string
		args    args
		want    media.Type
		wantErr bool
		err     error
	}{
		{
			name:    "found",
			args:    args{ctx: context.Background(), mediaId: m.Id},
			want:    media.TypeImageJPEG,
			wantErr: false,
			err:     nil,
		},
		{
			name:    "media not found",
			args:    args{ctx: context.Background(), mediaId: "4bc2892f-4883-41dd-a9f9-3010e7fbd131"},
			want:    "",
			wantErr: true,
			err:     media.ErrNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := _mediaRepo.FindMediaTypeById(tt.args.ctx, tt.args.mediaId)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.ErrorIs(t, err, tt.err)
		})
	}
}
