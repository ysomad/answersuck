package jwt

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicManager(t *testing.T) {
	t.Parallel()

	type args struct {
		sign   string
		issuer string
	}
	tests := []struct {
		name    string
		args    args
		want    basicManager
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sign:   "test_sign",
				issuer: "test_issuer",
			},
			want: basicManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			wantErr: false,
		},
		{
			name: "empty sign",
			args: args{
				sign: "",
			},
			want:    basicManager{},
			wantErr: true,
		},
		{
			name: "empty issuer",
			args: args{
				issuer: "",
			},
			want:    basicManager{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewBasicManager(tt.args.sign, tt.args.issuer)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSubOnlyManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSubOnlyManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_basicManager_Encode(t *testing.T) {
	t.Parallel()

	type args struct {
		claims BasicClaims
	}
	tests := []struct {
		name    string
		m       basicManager
		args    args
		wantErr bool
	}{
		{
			name: "success",
			m:    basicManager{sign: []byte("test_sign")},
			args: args{
				claims: BasicClaims{
					Subject:   "8551023e-ef42-4fee-87dc-76093c888125",
					ExpiresAt: time.Now().Add(time.Minute).Unix(),
					IssuedAt:  time.Now().Unix(),
					Issuer:    "test_issuer",
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.m.Encode(tt.args.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("subOnlyManager.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("subOnlyManager.Encode() = %v, want not empty string", got)
			}
		})
	}
}

func Test_basicManager_Parse(t *testing.T) {
	t.Parallel()

	type args struct {
		token Basic
	}
	tests := []struct {
		name    string
		m       basicManager
		args    args
		want    BasicClaims
		wantErr bool
	}{
		{
			name: "success",
			m: basicManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				token: Basic("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YiIsImlzcyI6InRlc3RfaXNzdWVyIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzUwMDB9.OCSyzmIPejQl7iiQD8ZLFdFw0-xx-UYwYJfC-YGQc_U"),
			},
			want: BasicClaims{
				ExpiresAt: 9223372036854775000,
				IssuedAt:  1516239022,
				Issuer:    "test_issuer",
				Subject:   "test_sub",
			},
			wantErr: false,
		},
		{
			name: "unexpected signing method",
			m: basicManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				// used PS512 for the test
				token: Basic("eyJhbGciOiJQUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YiIsImlzcyI6InRlc3RfaXNzdWVyIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzUwMDB9.MehhCh5glFg5I0y5iHkRTswq_ZiPpMuh5Pl1GfQ2UbrAUpCh8aUcDtN6uQmMb5Oz_5mV75igFmSHVAsjA8QwIXMWfL6-URXIVo8-N5yYv4gxoqpIsuT81vZ8JJXIG4U4hhDjHHAiPeykZZp7WlDaycgB4IGRwcsC4WMo1cqvx_dUCqUhLvCcziI4Wamn4GhYek9_q95LvAiK126-1YnxPiG3NmsOkzro8M6v052f_Y_LTDA85-gGlOPnZAI08jfA9myXDXe5wBvMLXMrw_Jv94mE2gEIA4Iaze9TzBHn6oASMHotVJ5bIeEdOCOj8gW6XImOPHjLHDRpjqfqqx580g"),
			},
			want:    BasicClaims{},
			wantErr: true,
		},
		{
			name: "invalid token",
			m: basicManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				// used PS512 for the test
				token: Basic("yeet"),
			},
			want:    BasicClaims{},
			wantErr: true,
		},
		{
			name: "invalid sign",
			m: basicManager{
				sign:   []byte("invalid_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				token: Basic("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YiIsImlzcyI6InRlc3RfaXNzdWVyIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzUwMDB9.OCSyzmIPejQl7iiQD8ZLFdFw0-xx-UYwYJfC-YGQc_U"),
			},
			want:    BasicClaims{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.m.Parse(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("subOnlyManager.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBasic_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tr   Basic
		want string
	}{
		{
			name: "success",
			tr:   Basic("yeet"),
			want: "yeet",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.tr.String(); got != tt.want {
				t.Errorf("SubOnly.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
