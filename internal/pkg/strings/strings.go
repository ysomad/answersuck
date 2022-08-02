package strings

import (
	cryptoRand "crypto/rand"
	"errors"
	mathRand "math/rand"
)

const (
	chars   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	special = " !\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

var errZeroLength = errors.New("length should be greater than 0")

// NewUnique generates random string using
// Cryptographically Secure Pseudorandom number
func NewUnique(length uint) (string, error) {
	if length == 0 {
		return "", errZeroLength
	}

	bytes := make([]byte, length)
	if _, err := cryptoRand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes), nil
}

// NewRandom generates random URL safe string
func NewRandom(length uint) string {
	bytes := make([]byte, length)

	for i := range bytes {
		bytes[i] = chars[mathRand.Intn(len(chars))]
	}

	return string(bytes)
}

// NewSpecialRandom generates random string with special characters
func NewSpecialRandom(length uint) string {
	bytes := make([]byte, length)

	c := chars + special
	for i := range bytes {
		bytes[i] = c[mathRand.Intn(len(chars))]
	}

	return string(bytes)
}
