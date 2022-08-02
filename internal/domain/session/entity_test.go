package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_newSession(t *testing.T) {
	t.Parallel()

	now := time.Now()
	type args struct {
		accountId string
		userAgent string
		ip        string
		exp       time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    *Session
		wantErr bool
	}{
		{
			name: "session created",
			args: args{
				accountId: "id",
				userAgent: "ua",
				ip:        "192.0.0.1",
				exp:       time.Minute,
			},
			want: &Session{
				AccountId: "id",
				UserAgent: "ua",
				IP:        "192.0.0.1",
				MaxAge:    int(time.Minute.Seconds()),
				ExpiresAt: now.Add(time.Minute).Unix(),
				CreatedAt: now,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newSession(tt.args.accountId, tt.args.userAgent, tt.args.ip, tt.args.exp)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.NotEqual(t, "", got.Id)
			assert.Equal(t, SessionIdLen, len(got.Id))
			assert.Equal(t, tt.want.AccountId, got.AccountId)
			assert.Equal(t, tt.want.UserAgent, got.UserAgent)
			assert.Equal(t, tt.want.IP, got.IP)
			assert.Equal(t, tt.want.MaxAge, got.MaxAge)
			assert.Equal(t, tt.want.ExpiresAt, got.ExpiresAt)
			assert.Equal(t, tt.want.CreatedAt.Unix(), got.CreatedAt.Unix())
		})
	}
}

func TestSession_Expired(t *testing.T) {
	type fields struct {
		ExpiresAt int64
		CreatedAt time.Time
	}
	now := time.Now()
	tests := []struct {
		name      string
		expiresAt int64
		want      bool
	}{
		{
			name:      "session expired",
			expiresAt: now.Add(-time.Hour).Unix(),
			want:      true,
		},
		{
			name:      "session not expired",
			expiresAt: now.Add(time.Hour).Unix(),
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{ExpiresAt: tt.expiresAt}
			got := s.Expired()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSession_SameDevice(t *testing.T) {
	type fields struct {
		UserAgent string
		IP        string
	}
	tests := []struct {
		name   string
		fields fields
		device Device
		want   bool
	}{
		{
			name: "same device",
			fields: fields{
				UserAgent: "linux chrome",
				IP:        "133.7.42.0",
			},
			device: Device{
				UserAgent: "linux chrome",
				IP:        "133.7.42.0",
			},
			want: true,
		},
		{
			name: "different user agent",
			fields: fields{
				UserAgent: "firefox macos",
				IP:        "133.7.42.0",
			},
			device: Device{
				UserAgent: "ua",
				IP:        "133.7.42.0",
			},
			want: false,
		},
		{
			name: "different ip",
			fields: fields{
				UserAgent: "firefox macos",
				IP:        "133.7.42.0",
			},
			device: Device{
				UserAgent: "firefox macos",
				IP:        "148.8.133.7",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				UserAgent: tt.fields.UserAgent,
				IP:        tt.fields.IP,
			}
			got := s.SameDevice(Device{IP: tt.device.IP, UserAgent: tt.device.UserAgent})
			assert.Equal(t, tt.want, got)
		})
	}
}
