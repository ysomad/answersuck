package account

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPasswordToken_expired(t *testing.T) {
	t.Parallel()

	type args struct {
		exp time.Duration
	}
	tests := []struct {
		name      string
		createdAt time.Time
		args      args
		want      bool
	}{
		{
			name:      "password token expired",
			createdAt: time.Now().Add(-time.Hour),
			args:      args{exp: time.Minute},
			want:      true,
		},
		{
			name:      "password token not expired",
			createdAt: time.Now(),
			args:      args{exp: time.Minute},
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := PasswordToken{CreatedAt: tt.createdAt}
			got := tr.expired(tt.args.exp)
			assert.Equal(t, tt.want, got)
		})
	}
}
