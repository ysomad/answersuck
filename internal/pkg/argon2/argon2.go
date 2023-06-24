package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/crypto/argon2"
)

var (
	errInvalidHash         = errors.New("argon2: hash is not in the correct format")
	errIncompatibleVariant = errors.New("argon2: incompatible variant of argon2")
	errIncompatibleVersion = errors.New("argon2: incompatible version of argon2")
	errInvalidSaltLength   = errors.New("argon2: invalid params salt length")
	errEmptyPassword       = errors.New("argon2: empty password")
)

var DefaultParams = Params{ //nolint:gochecknoglobals // it has to be here to use the package
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

type Params struct {
	// The amount of Memory used by the algorithm (in kibibytes).
	Memory uint32

	// The number of Iterations over the memory.
	Iterations uint32

	// The number of threads (or lanes) used by the algorithm.
	// Recommended value is between 1 and runtime.NumCPU().
	Parallelism uint8

	// Length of the random salt. 16 bytes is recommended for password hashing.
	SaltLength uint32

	// Length of the generated key. 16 bytes or more is recommended.
	KeyLength uint32
}

func GenerateFromPassword(password string, p Params) (hash string, err error) {
	if password == "" {
		return "", errEmptyPassword
	}

	if p.SaltLength <= 0 {
		return "", errInvalidSaltLength
	}

	salt := make([]byte, p.SaltLength)

	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := argon2.IDKey(
		unsafe.Slice(unsafe.StringData(password), len(password)),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		p.Memory,
		p.Iterations,
		p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func CompareHashAndPassword(password, hash string) (bool, error) {
	p, salt, key, err := decode(hash)
	if err != nil {
		return false, err
	}

	comparedKey := argon2.IDKey(
		unsafe.Slice(unsafe.StringData(password), len(password)),
		salt,
		p.Iterations,
		p.Memory,
		p.Parallelism,
		p.KeyLength,
	)

	if subtle.ConstantTimeEq(int32(len(key)), int32(len(comparedKey))) == 0 {
		return false, nil
	}

	if subtle.ConstantTimeCompare(key, comparedKey) == 1 {
		return true, nil
	}

	return false, nil
}

// decode expects a hash created from this package, and parses it to return the params used to
// create it, as well as the salt and key (password hash).
func decode(hash string) (params Params, salt, key []byte, err error) {
	vals := strings.Split(hash, "$")
	if len(vals) != 6 {
		return Params{}, nil, nil, errInvalidHash
	}

	if vals[1] != "argon2id" {
		return Params{}, nil, nil, errIncompatibleVariant
	}

	var version int

	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return Params{}, nil, nil, err
	}

	if version != argon2.Version {
		return Params{}, nil, nil, errIncompatibleVersion
	}

	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &params.Memory, &params.Iterations, &params.Parallelism)
	if err != nil {
		return Params{}, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return Params{}, nil, nil, err
	}

	params.SaltLength = uint32(len(salt))

	key, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return Params{}, nil, nil, err
	}

	params.KeyLength = uint32(len(key))

	return params, salt, key, nil
}
