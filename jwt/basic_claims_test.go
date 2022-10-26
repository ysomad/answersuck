package jwt

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNewBasicClaims(t *testing.T) {
	t.Parallel()

	type args struct {
		subject string
		issuer  string
		exp     time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    BasicClaims
		wantErr error
	}{
		{
			name: "success",
			args: args{
				subject: "8551023e-ef42-4fee-87dc-76093c888125",
				issuer:  "test_issuer",
				exp:     time.Hour,
			},
			want: BasicClaims{
				Subject:   "8551023e-ef42-4fee-87dc-76093c888125",
				Issuer:    "test_issuer",
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Hour).Unix(),
			},
			wantErr: nil,
		},
		{
			name: "invalid subject",
			args: args{
				subject: "invalid_sub",
				issuer:  "test_issuer",
				exp:     time.Hour,
			},
			want:    BasicClaims{},
			wantErr: errInvalidSubject,
		},
		{
			name: "empty subject",
			args: args{
				subject: "",
				issuer:  "test_issuer",
				exp:     time.Hour,
			},
			want:    BasicClaims{},
			wantErr: errInvalidSubject,
		},
		{
			name: "empty issuer",
			args: args{
				subject: "8551023e-ef42-4fee-87dc-76093c888125",
				issuer:  "",
				exp:     time.Hour,
			},
			want:    BasicClaims{},
			wantErr: errInvalidIssuer,
		},
		{
			name: "0 exp",
			args: args{
				subject: "8551023e-ef42-4fee-87dc-76093c888125",
				issuer:  "test_issuer",
				exp:     0,
			},
			want:    BasicClaims{},
			wantErr: errInvalidExpiration,
		},
		{
			name: "negative exp",
			args: args{
				subject: "8551023e-ef42-4fee-87dc-76093c888125",
				issuer:  "test_issuer",
				exp:     -1,
			},
			want:    BasicClaims{},
			wantErr: errInvalidExpiration,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewBasicClaims(tt.args.subject, tt.args.issuer, tt.args.exp)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBasicClaims_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		c       BasicClaims
		wantErr bool
	}{
		{
			name: "valid",
			c: BasicClaims{
				Subject:   "8551023e-ef42-4fee-87dc-76093c888125",
				ExpiresAt: time.Now().Add(time.Minute).Unix(),
			},
			wantErr: false,
		},
		{
			name: "invalid Subject, valid expiresAt",
			c: BasicClaims{
				Subject:   "invalid_subject",
				ExpiresAt: time.Now().Add(time.Minute).Unix(),
			},
			wantErr: true,
		},
		{
			name: "empty Subject, valid expiresAt",
			c: BasicClaims{
				Subject:   "",
				ExpiresAt: time.Now().Add(time.Minute).Unix(),
			},
			wantErr: true,
		},
		{
			name: "empty Subject, empty ExpiresAt",
			c: BasicClaims{
				Subject:   "",
				ExpiresAt: 0,
			},
			wantErr: true,
		},
		{
			name: "empty Subject, negative ExpiresAt",
			c: BasicClaims{
				Subject:   "",
				ExpiresAt: -1,
			},
			wantErr: true,
		},
		{
			name: "valid Subject, negative ExpiresAt",
			c: BasicClaims{
				Subject:   "8551023e-ef42-4fee-87dc-76093c888125",
				ExpiresAt: -1,
			},
			wantErr: true,
		},
		{
			name: "valid Subject, empty ExpiresAt",
			c: BasicClaims{
				Subject:   "8551023e-ef42-4fee-87dc-76093c888125",
				ExpiresAt: 0,
			},
			wantErr: true,
		},
		{
			name: "expired",
			c: BasicClaims{
				Subject:   "8551023e-ef42-4fee-87dc-76093c888125",
				ExpiresAt: time.Now().Add(-time.Minute).Unix(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := tt.c.Valid(); (err != nil) != tt.wantErr {
				t.Errorf("SubOnlyClaims.Valid() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBasicClaims_verifyExpiresAt(t *testing.T) {
	t.Parallel()

	type args struct {
		now int64
	}
	tests := []struct {
		name string
		c    BasicClaims
		args args
		want bool
	}{
		{
			name: "expired",
			c: BasicClaims{
				ExpiresAt: time.Now().Add(-time.Hour).Unix(),
			},
			args: args{
				now: time.Now().Unix(),
			},
			want: false,
		},
		{
			name: "empty ExpiresAt",
			c:    BasicClaims{},
			args: args{
				now: time.Now().Unix(),
			},
			want: false,
		},
		{
			name: "verified",
			c: BasicClaims{
				ExpiresAt: time.Now().Add(time.Minute).Unix(),
			},
			args: args{
				now: time.Now().Unix(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.c.verifyExpiresAt(tt.args.now); got != tt.want {
				t.Errorf("SubOnlyClaims.verifyExpiresAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicClaims_verifySubject(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    BasicClaims
		want bool
	}{
		{
			name: "valid subject",
			c: BasicClaims{
				Subject: "8551023e-ef42-4fee-87dc-76093c888125",
			},
			want: true,
		},
		{
			name: "invalid subject",
			c: BasicClaims{
				Subject: "invalid_subject",
			},
			want: false,
		},
		{
			name: "empty subject",
			c: BasicClaims{
				Subject: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.c.verifySubject(); got != tt.want {
				t.Errorf("SubOnlyClaims.verifySubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicClaims_verifyIssuer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		c    BasicClaims
		want bool
	}{
		{
			name: "verified",
			c: BasicClaims{
				Issuer: "test_issuer",
			},
			want: true,
		},
		{
			name: "invalid issuer",
			c: BasicClaims{
				Issuer: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.c.verifyIssuer(); got != tt.want {
				t.Errorf("SubOnlyClaims.verifyIssuer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBasicClaims_verifyIssuedAt(t *testing.T) {
	t.Parallel()

	type args struct {
		now int64
	}
	tests := []struct {
		name string
		c    BasicClaims
		args args
		want bool
	}{
		{
			name: "verified",
			c: BasicClaims{
				IssuedAt: time.Now().Add(-time.Minute).Unix(),
			},
			args: args{
				now: time.Now().Unix(),
			},
			want: true,
		},
		{
			name: "used before issued",
			c: BasicClaims{
				IssuedAt: time.Now().Add(time.Minute).Unix(),
			},
			args: args{
				now: time.Now().Unix(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.c.verifyIssuedAt(tt.args.now); got != tt.want {
				t.Errorf("SubOnlyClaims.verifyIssuedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newBasicClaims(t *testing.T) {
	t.Parallel()

	type args struct {
		raw jwt.MapClaims
	}
	tests := []struct {
		name    string
		args    args
		want    BasicClaims
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				jwt.MapClaims{
					"exp": time.Now().Add(time.Minute).Unix(),
					"iat": time.Now().Add(-time.Minute).Unix(),
					"sub": "test_subject",
					"iss": "test_issuer",
				},
			},
			want: BasicClaims{
				ExpiresAt: time.Now().Add(time.Minute).Unix(),
				IssuedAt:  time.Now().Add(-time.Minute).Unix(),
				Subject:   "test_subject",
				Issuer:    "test_issuer",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := newBasicClaims(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("newSubOnlyClaims() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSubOnlyClaims() = %v, want %v", got, tt.want)
			}
		})
	}
}
