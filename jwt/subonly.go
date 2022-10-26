package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// SubOnly is the json web token contains only subject.
type SubOnly string

type subOnlyManager struct {
	sign   []byte
	Issuer string
}

func NewSubOnlyManager(sign, issuer string) (subOnlyManager, error) {
	if sign == "" {
		return subOnlyManager{}, errEmptySign
	}
	if issuer == "" {
		return subOnlyManager{}, errInvalidIssuer
	}

	return subOnlyManager{sign: []byte(sign)}, nil
}

func (m subOnlyManager) Encode(claims SubOnlyClaims) (SubOnly, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(m.sign)
	if err != nil {
		return "", err
	}

	return SubOnly(s), nil
}

func (m subOnlyManager) Parse(token SubOnly) (SubOnlyClaims, error) {
	t, err := jwt.Parse(string(token), func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSignMethod
		}

		return m.sign, nil
	})
	if err != nil {
		return SubOnlyClaims{}, err
	}

	raw, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return SubOnlyClaims{}, errNoClaims
	}

	c, err := newSubOnlyClaims(raw)
	if err != nil {
		return SubOnlyClaims{}, fmt.Errorf("%s: %w", err.Error(), errNoClaims)
	}

	return c, nil
}
