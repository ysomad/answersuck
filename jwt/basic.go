package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt"
)

// basicManager is the manager for token with BasicClaims.
type basicManager struct {
	sign   []byte
	issuer string
}

func NewBasicManager(sign, issuer string) (basicManager, error) {
	if sign == "" {
		return basicManager{}, errEmptySign
	}
	if issuer == "" {
		return basicManager{}, errInvalidIssuer
	}

	return basicManager{
		sign:   []byte(sign),
		issuer: issuer,
	}, nil
}

func (m basicManager) Encode(claims BasicClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	s, err := t.SignedString(m.sign)
	if err != nil {
		return "", err
	}

	return s, nil
}

func (m basicManager) Decode(token string) (BasicClaims, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
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

func (m basicManager) Issuer() string { return m.issuer }
