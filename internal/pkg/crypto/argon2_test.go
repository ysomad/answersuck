package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_argon2Id(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		hash    string
		ok      bool
		wantErr bool
	}{
		{
			name:    "correct match",
			plain:   "inanzzz",
			hash:    "$argon2id$v=19$m=65536,t=1,p=4$3Ymrg20KqXd9mMawXP/YzA$nf7ubeO0tB1NDk4nBscgsHvIcDECMjIuEeEjgBRMe3s",
			ok:      true,
			wantErr: false,
		},
		{
			name:    "incorrect match",
			plain:   "test123",
			hash:    "$argon2id$v=19$m=65536,t=1,p=4$L/6i/DD9Ie5dKo7L6PpvVg$8OLc5G2E715nbsSgA4ZKcGngVLtAeCnB4CD76XbShic",
			ok:      false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArgon2Id()
			got, err := a.Verify(tt.plain, tt.hash)

			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.ok, got)
		})
	}
}
