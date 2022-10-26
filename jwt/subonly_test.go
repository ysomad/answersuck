package jwt

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSubOnlyManager(t *testing.T) {
	t.Parallel()

	type args struct {
		sign   string
		issuer string
	}
	tests := []struct {
		name    string
		args    args
		want    subOnlyManager
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				sign:   "test_sign",
				issuer: "test_issuer",
			},
			want: subOnlyManager{
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
			want:    subOnlyManager{},
			wantErr: true,
		},
		{
			name: "empty issuer",
			args: args{
				issuer: "",
			},
			want:    subOnlyManager{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewSubOnlyManager(tt.args.sign, tt.args.issuer)
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

func Test_subOnlyManager_Encode(t *testing.T) {
	t.Parallel()

	type args struct {
		claims SubOnlyClaims
	}
	tests := []struct {
		name    string
		m       subOnlyManager
		args    args
		wantErr bool
	}{
		{
			name: "success",
			m:    subOnlyManager{sign: []byte("test_sign")},
			args: args{
				claims: SubOnlyClaims{
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

func Test_subOnlyManager_Parse(t *testing.T) {
	t.Parallel()

	type args struct {
		token SubOnly
	}
	tests := []struct {
		name    string
		m       subOnlyManager
		args    args
		want    SubOnlyClaims
		wantErr bool
	}{
		{
			name: "success",
			m: subOnlyManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				token: SubOnly("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YiIsImlzcyI6InRlc3RfaXNzdWVyIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzUwMDB9.OCSyzmIPejQl7iiQD8ZLFdFw0-xx-UYwYJfC-YGQc_U"),
			},
			want: SubOnlyClaims{
				ExpiresAt: 9223372036854775000,
				IssuedAt:  1516239022,
				Issuer:    "test_issuer",
				Subject:   "test_sub",
			},
			wantErr: false,
		},
		{
			name: "unexpected signing method",
			m: subOnlyManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				// used PS512 for the test
				token: SubOnly("eyJhbGciOiJQUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YiIsImlzcyI6InRlc3RfaXNzdWVyIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzUwMDB9.MehhCh5glFg5I0y5iHkRTswq_ZiPpMuh5Pl1GfQ2UbrAUpCh8aUcDtN6uQmMb5Oz_5mV75igFmSHVAsjA8QwIXMWfL6-URXIVo8-N5yYv4gxoqpIsuT81vZ8JJXIG4U4hhDjHHAiPeykZZp7WlDaycgB4IGRwcsC4WMo1cqvx_dUCqUhLvCcziI4Wamn4GhYek9_q95LvAiK126-1YnxPiG3NmsOkzro8M6v052f_Y_LTDA85-gGlOPnZAI08jfA9myXDXe5wBvMLXMrw_Jv94mE2gEIA4Iaze9TzBHn6oASMHotVJ5bIeEdOCOj8gW6XImOPHjLHDRpjqfqqx580g"),
			},
			want:    SubOnlyClaims{},
			wantErr: true,
		},
		{
			name: "invalid token",
			m: subOnlyManager{
				sign:   []byte("test_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				// used PS512 for the test
				token: SubOnly("yeet"),
			},
			want:    SubOnlyClaims{},
			wantErr: true,
		},
		{
			name: "invalid sign",
			m: subOnlyManager{
				sign:   []byte("invalid_sign"),
				Issuer: "test_issuer",
			},
			args: args{
				token: SubOnly("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0X3N1YiIsImlzcyI6InRlc3RfaXNzdWVyIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjkyMjMzNzIwMzY4NTQ3NzUwMDB9.OCSyzmIPejQl7iiQD8ZLFdFw0-xx-UYwYJfC-YGQc_U"),
			},
			want:    SubOnlyClaims{},
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
