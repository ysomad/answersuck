package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	errEmptySign            = errors.New("signing key is empty")
	errNoClaims             = errors.New("error getting claims from token")
	errUnexpectedSignMethod = errors.New("unexpected signing method")
)

type manager struct {
	sign string
}

func NewManager(sign string) (manager, error) {
	if sign == "" {
		return manager{}, errEmptySign
	}

	return manager{
		sign: sign,
	}, nil
}

func (tm manager) Create(subject string, expiration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   subject,
		ExpiresAt: time.Now().Add(expiration).Unix(),
	})

	return token.SignedString([]byte(tm.sign))
}

func (tm manager) Parse(token string) (string, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (i interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSignMethod
		}

		return []byte(tm.sign), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return "", errNoClaims
	}

	return claims["sub"].(string), nil
}
