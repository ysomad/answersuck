package argon2

import (
	"testing"
)

func TestCompare(t *testing.T) {
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
		{
			name:    "incorrect hash",
			plain:   "test123",
			hash:    "invalid_hash",
			ok:      false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			argon2 := New()
			got, err := argon2.Compare(tt.plain, tt.hash)

			if (err != nil) != tt.wantErr {
				t.Errorf("want wantErr %t, got %s", tt.wantErr, err.Error())
			}

			if tt.ok != got {
				t.Errorf("want %t, got %t", tt.ok, got)
			}
		})
	}
}
