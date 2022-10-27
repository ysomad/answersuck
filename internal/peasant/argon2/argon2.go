package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2ID struct {
	version int
	params  *params
}

type params struct {
	memory      uint32
	iterations  uint32
	saltLen     uint32
	keyLen      uint32
	parallelism uint8
}

var DefaultParams = &params{
	memory:      64 * 1024,
	iterations:  1,
	parallelism: 2,
	saltLen:     16,
	keyLen:      32,
}

func NewID(p *params) *argon2ID {
	return &argon2ID{
		version: argon2.Version,
		params:  p,
	}
}

// Encode encodes plain text to argon2 hash
func (a *argon2ID) Encode(plain string) (string, error) {
	salt := make([]byte, a.params.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(plain), salt, a.params.iterations, a.params.memory, a.params.parallelism, a.params.keyLen)
	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		a.version,
		a.params.memory,
		a.params.iterations,
		a.params.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

// Compare compares plain text with hash and returns true, nil if it matches
func (a *argon2ID) Compare(plain, encoded string) (bool, error) {
	p, salt, key, err := decodeID(encoded)
	if err != nil {
		return false, err
	}

	keyToCompare := argon2.IDKey([]byte(plain), salt, p.iterations, p.memory, p.parallelism, p.keyLen)

	if subtle.ConstantTimeEq(int32(len(key)), int32(len(keyToCompare))) == 0 {
		return false, nil
	}

	return subtle.ConstantTimeCompare(key, keyToCompare) == 1, nil
}

func decodeID(hash string) (p *params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, errors.New("invalid hash")
	}

	if vals[1] != "argon2id" {
		return nil, nil, nil, errors.New("incompatible variant")
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errors.New("incompatible version")
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLen = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLen = uint32(len(key))

	return p, salt, key, nil
}
