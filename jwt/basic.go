package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// Basic is the json web token contains only sub, iat, iss and exp.
type Basic string

func (b Basic) String() string { return string(b) }

type basicManager struct {
	sign   []byte
	Issuer string
}

func NewBasicManager(sign, issuer string) (basicManager, error) {
	if sign == "" {
		return basicManager{}, errEmptySign
	}
	if issuer == "" {
		return basicManager{}, errInvalidIssuer
	}

	return basicManager{sign: []byte(sign)}, nil
}

func (m basicManager) Encode(claims BasicClaims) (Basic, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(m.sign)
	if err != nil {
		return "", err
	}

	return Basic(s), nil
}

func (m basicManager) Parse(token Basic) (BasicClaims, error) {
	t, err := jwt.Parse(string(token), func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSignMethod
		}

		return m.sign, nil
	})
	if err != nil {
		return BasicClaims{}, err
	}

	raw, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return BasicClaims{}, errNoClaims
	}

	c, err := newBasicClaims(raw)
	if err != nil {
		return BasicClaims{}, fmt.Errorf("%s: %w", err.Error(), errNoClaims)
	}

	return c, nil
}
