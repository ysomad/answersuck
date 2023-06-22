package argon2

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/argon2"
)

// testHash is a test argon2id hash with DefaultParams with key "testtest".
const testHash = "$argon2id$v=19$m=65536,t=1,p=2$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU"

func TestGenerateFromPassword(t *testing.T) {
	// invalid salt length
	got1, err := GenerateFromPassword("test", Params{SaltLength: 0})
	assert.Error(t, err)
	assert.Equal(t, "", got1)

	// success
	got2, err := GenerateFromPassword("test", DefaultParams)
	assert.NoError(t, err)

	vals := strings.Split(got2, "$")
	assert.Equal(t, 6, len(vals))
	assert.Equal(t, "argon2id", vals[1])

	var (
		ver int
		p   Params
	)

	_, err = fmt.Sscanf(vals[2], "v=%d", &ver)
	assert.NoError(t, err)
	assert.Equal(t, argon2.Version, ver)

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	assert.NoError(t, err)

	salt, err := base64.RawStdEncoding.Strict().DecodeString(vals[4])
	assert.NoError(t, err)
	assert.Equal(t, int(DefaultParams.SaltLength), len(salt))

	key, err := base64.RawStdEncoding.Strict().DecodeString(vals[5])
	assert.NoError(t, err)
	assert.Equal(t, int(DefaultParams.KeyLength), len(key))

	// empty password
	got3, err := GenerateFromPassword("", DefaultParams)
	assert.Error(t, err)
	assert.Equal(t, "", got3)
}

func Test_decode(t *testing.T) {
	t.Parallel()

	type args struct {
		hash string
	}

	testSalt, _ := base64.RawStdEncoding.Strict().DecodeString("mbmLJMzdBa1rGVJGCT+mxA")
	testKey, _ := base64.RawStdEncoding.Strict().DecodeString("skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU")

	tests := []struct {
		name       string
		args       args
		wantParams Params
		wantSalt   []byte
		wantKey    []byte
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				hash: testHash,
			},
			wantParams: DefaultParams,
			wantSalt:   testSalt,
			wantKey:    testKey,
			wantErr:    false,
		},
		{
			name: "invalid hash length",
			args: args{
				hash: "$argon2id$v=19$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU", // no params
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "invalid hash prefix",
			args: args{
				hash: "$INVALID$v=19$m=65536,t=1,p=2$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "invalid version",
			args: args{
				hash: "$argon2id$INVALID=19$m=65536,t=1,p=2$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "incompatible version",
			args: args{
				hash: "$argon2id$v=18$m=65536,t=1,p=2$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "invalid memory",
			args: args{
				hash: "$argon2id$v=19$m=INVALID,t=1,p=2$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "invalid iterations",
			args: args{
				hash: "$argon2id$v=19$m=65536,t=INVALID,p=2$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "invalid parallelism",
			args: args{
				hash: "$argon2id$v=19$m=65536,t=1,p=INVALID$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "corrupted salt",
			args: args{
				hash: "$argon2id$v=19$m=65536,t=1,p=2$XXXXXaGVsbG8=$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
		{
			name: "corrupted key",
			args: args{
				hash: "$argon2id$v=19$m=65536,t=1,p=2$mbmLJMzdBa1rGVJGCT+mxA$XXXXXaGVsbG8=",
			},
			wantParams: Params{},
			wantSalt:   nil,
			wantKey:    nil,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParams, gotSalt, gotKey, err := decode(tt.args.hash)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.wantParams, gotParams)
			assert.Equal(t, tt.wantSalt, gotSalt)
			assert.Equal(t, tt.wantKey, gotKey)
		})
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	type args struct {
		password string
		hash     string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "decode error",
			args: args{
				password: "testtest",
				hash:     "$argon2id$v=19$mbmLJMzdBa1rGVJGCT+mxA$skLbmhX8uevw262bn0i5wREOY42a/1Sm588IZnceBPU",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid key",
			args: args{
				password: "invalid",
				hash:     testHash,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "valid key",
			args: args{
				password: "testtest",
				hash:     testHash,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompareHashAndPassword(tt.args.password, tt.args.hash)
			assert.Equal(t, tt.wantErr, (err != nil))
			assert.Equal(t, tt.want, got)
		})
	}
}
