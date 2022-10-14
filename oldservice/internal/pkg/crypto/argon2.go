package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2Id struct {
	format  string
	version int
	time    uint32
	memory  uint32
	keyLen  uint32
	saltLen uint32
	threads uint8
}

func NewArgon2Id() *argon2Id {
	return &argon2Id{
		format:  "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		version: argon2.Version,
		time:    1,
		memory:  64 * 1024,
		keyLen:  32,
		saltLen: 32,
		threads: 4,
	}
}

func (a argon2Id) Hash(plain string) (string, error) {
	salt := make([]byte, a.saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(plain), salt, a.time, a.memory, a.threads, a.keyLen)
	return fmt.Sprintf(
		a.format,
		a.version,
		a.memory,
		a.time,
		a.threads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func (a argon2Id) Verify(plain, hash string) (bool, error) {
	hashParts := strings.Split(hash, "$")

	_, err := fmt.Sscanf(hashParts[3], "m=%d,t=%d,p=%d", &a.memory, &a.time, &a.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(hashParts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(hashParts[5])
	if err != nil {
		return false, err
	}

	hashToCompare := argon2.IDKey([]byte(plain), salt, a.time, a.memory, a.threads, uint32(len(decodedHash)))

	return subtle.ConstantTimeCompare(decodedHash, hashToCompare) == 1, nil
}
