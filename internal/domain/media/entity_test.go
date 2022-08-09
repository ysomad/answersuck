package media

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestType_valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mt   Type
		want bool
	}{
		{
			name: "valid jpeg",
			mt:   Type("image/jpeg"),
			want: true,
		},
		{
			name: "valid mp4",
			mt:   Type("audio/mp4"),
			want: true,
		},
		{
			name: "valid png",
			mt:   Type("image/png"),
			want: true,
		},
		{
			name: "valid aac",
			mt:   Type("audio/aac"),
			want: true,
		},
		{
			name: "valid mpeg",
			mt:   Type("audio/mpeg"),
			want: true,
		},
		{
			name: "invalid media type",
			mt:   Type("invalid"),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.mt.valid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		filename  string
		accountId string
		t         Type
	}
	tests := []struct {
		name    string
		args    args
		want    Media
		wantErr bool
		err     error
	}{
		{
			name: "valid media png",
			args: args{
				filename:  "test",
				accountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				t:         Type("image/png"),
			},
			want: Media{
				Id:        "some generated id",
				Filename:  "test",
				Type:      Type("image/png"),
				AccountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				CreatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "valid media jpeg",
			args: args{
				filename:  "test",
				accountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				t:         Type("image/jpeg"),
			},
			want: Media{
				Id:        "some generated id",
				Filename:  "test",
				Type:      Type("image/jpeg"),
				AccountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				CreatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "valid media mp4",
			args: args{
				filename:  "test",
				accountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				t:         Type("audio/mp4"),
			},
			want: Media{
				Id:        "some generated id",
				Filename:  "test",
				Type:      Type("audio/mp4"),
				AccountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				CreatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "valid media aac",
			args: args{
				filename:  "test",
				accountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				t:         Type("audio/aac"),
			},
			want: Media{
				Id:        "some generated id",
				Filename:  "test",
				Type:      Type("audio/aac"),
				AccountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				CreatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "valid media mpeg",
			args: args{
				filename:  "test",
				accountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				t:         Type("audio/mpeg"),
			},
			want: Media{
				Id:        "some generated id",
				Filename:  "test",
				Type:      Type("audio/mpeg"),
				AccountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				CreatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "valid media image webp",
			args: args{
				filename:  "test",
				accountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				t:         Type("image/webp"),
			},
			want: Media{
				Id:        "some generated id",
				Filename:  "test",
				Type:      Type("image/webp"),
				AccountId: "58f0eb78-5080-46ee-8a6d-18950477bba0",
				CreatedAt: time.Now(),
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.filename, tt.args.accountId, tt.args.t)
			assert.ErrorIs(t, err, tt.err)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, (tt.want.Id != ""), (got.Id != ""))
			assert.Contains(t, got.Filename, tt.want.Filename)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.AccountId, got.AccountId)
			assert.Equal(t, false, got.CreatedAt.IsZero())
		})
	}
}
