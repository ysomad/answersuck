package media

import (
	"reflect"
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.filename, tt.args.accountId, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMedia_removeTmpFile(t *testing.T) {
	type fields struct {
		Id        string
		Filename  string
		Type      Type
		AccountId string
		CreatedAt time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Media{
				Id:        tt.fields.Id,
				Filename:  tt.fields.Filename,
				Type:      tt.fields.Type,
				AccountId: tt.fields.AccountId,
				CreatedAt: tt.fields.CreatedAt,
			}
			m.removeTmpFile()
		})
	}
}
