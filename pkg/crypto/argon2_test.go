package crypto

import (
	"testing"
)

func Test_argon2Id(t *testing.T) {
	tests := []struct {
		name    string
		plain   string
		hash    string
		want    bool
		wantErr bool
	}{
		{
			name:    "correct match",
			plain:   "inanzzz",
			hash:    "$argon2id$v=19$m=65536,t=1,p=4$3Ymrg20KqXd9mMawXP/YzA$nf7ubeO0tB1NDk4nBscgsHvIcDECMjIuEeEjgBRMe3s",
			want:    true,
			wantErr: false,
		},
		{
			name:    "incorrect match",
			plain:   "test123",
			hash:    "$argon2id$v=19$m=65536,t=1,p=4$L/6i/DD9Ie5dKo7L6PpvVg$8OLc5G2E715nbsSgA4ZKcGngVLtAeCnB4CD76XbShic",
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := NewArgon2Id()
			got, err := a.Verify(tt.plain, tt.hash)

			if (err != nil) != tt.wantErr {
				t.Errorf("argon2Id.Match() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("argon2Id.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
