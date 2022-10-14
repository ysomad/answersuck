// Package cryptostr provides functions for generating random strings using cryptographically
// secure random number generator from `crypto/rand`.
package cryptostr

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// RandomBase64 returns URL-safe, base64 encoded random string.
// For example It can be used for session id or csrf tokens.
func RandomBase64(n int) (string, error) {
	b, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// RandomWithCharset returns securely generated random string from given charset.
func RandomWithCharset(n int, cs charset) (string, error) {
	b := make([]byte, n)

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(cs))))
		if err != nil {
			return "", err
		}

		b[i] = cs[num.Int64()]
	}

	return string(b), nil
}

func Random(n int) (string, error)                   { return RandomWithCharset(n, AllChars) }
func RandomWithAlphabetLower(n int) (string, error)  { return RandomWithCharset(n, AlphabetLower) }
func RandomWithAlphabetUpper(n int) (string, error)  { return RandomWithCharset(n, AlphabetUpper) }
func RandomWithAlphabet(n int) (string, error)       { return RandomWithCharset(n, Alphabet) }
func RandomWithSpecials(n int) (string, error)       { return RandomWithCharset(n, Specials) }
func RandomWithAlphabetDigits(n int) (string, error) { return RandomWithCharset(n, AlphabetDigits) }
func RandomWithDigits(n int) (string, error)         { return RandomWithCharset(n, Digits) }

// RandomWithCharsets returns random string generated from slice of charsets,
// if sets are not specified, string will be generated from alphabet
// with lower, uppercase characters and digits.
//
// Better to use `NewWithCharset` or `NewWith{{charset}}` functions for better perfomance
// or `NewRandom` if all characters are needed.
func RandomWithCharsets(n int, sets ...charset) (string, error) {
	var chars charset
	if len(sets) > 0 {
		for _, set := range sets {
			chars += set
		}
	} else {
		chars = AlphabetDigits
	}

	return RandomWithCharset(n, charset(chars))
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
