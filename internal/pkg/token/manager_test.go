package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func createTestToken(sub, sign string, expiresAt int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   sub,
		ExpiresAt: expiresAt,
	})
	return token.SignedString([]byte(sign))
}

func TestNewManager(t *testing.T) {
	t.Parallel()

	type args struct {
		sign string
	}
	tests := []struct {
		name    string
		args    args
		want    manager
		wantErr bool
		err     error
	}{
		{
			name:    "success",
			args:    args{sign: "sign"},
			want:    manager{sign: "sign"},
			wantErr: false,
			err:     nil,
		},
		{
			name:    "empty sign",
			args:    args{sign: ""},
			want:    manager{},
			wantErr: true,
			err:     errEmptySign,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManager(tt.args.sign)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				assert.Equal(t, tt.want, got)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_manager_Create(t *testing.T) {
	t.Parallel()

	type fields struct {
		sign string
	}
	type args struct {
		subject    string
		expiration time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		err     error
	}{
		{
			name:    "token created",
			fields:  fields{sign: "sign"},
			args:    args{subject: "accountId", expiration: time.Minute},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := manager{
				sign: tt.fields.sign,
			}
			got, err := tm.Create(tt.args.subject, tt.args.expiration)
			if tt.wantErr {
				assert.ErrorIs(t, err, tt.err)
				return
			}
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NoError(t, err)

			assert.NotEqual(t, "", got)
		})
	}
}

func Test_manager_Parse(t *testing.T) {
	t.Parallel()

	type fields struct {
		sign string
	}
	type args struct {
		token string
	}

	sub := "sub"
	sign := "sign"
	token, err := createTestToken(sub, sign, time.Now().Add(time.Minute).Unix())
	assert.NoError(t, err)

	expiredToken, err := createTestToken(sub, sign, time.Now().Add(-time.Minute).Unix())
	assert.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{
			name:    "token is valid",
			fields:  fields{sign: sign},
			args:    args{token: token},
			want:    sub,
			wantErr: false,
		},
		{
			name:    "token is expired",
			fields:  fields{sign: sign},
			args:    args{token: expiredToken},
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid sign",
			fields:  fields{sign: ""},
			args:    args{token: token},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := manager{sign: tt.fields.sign}
			got, err := tm.Parse(tt.args.token)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}
