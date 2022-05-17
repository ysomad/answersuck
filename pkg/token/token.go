package token

import (
	"errors"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrEmptySign            = errors.New("signing key is empty")
	ErrNoClaims             = errors.New("error getting claims from token")
	ErrUnexpectedSignMethod = errors.New("unexpected signing method")
	ErrInvalidAudience      = errors.New("audience is not valid uri")
	ErrAudienceMismatch     = errors.New("audience does not match the audience in the token")
)

type manager struct {
	sign string
}

func NewManager(sign string) (manager, error) {
	if sign == "" {
		return manager{}, ErrEmptySign
	}

	return manager{
		sign: sign,
	}, nil
}

func (tm manager) New(subject, audience string, expiration time.Duration) (string, error) {
	_, err := url.Parse(audience)
	if err != nil {
		return "", ErrInvalidAudience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   subject,
		Audience:  audience,
		ExpiresAt: time.Now().Add(expiration).Unix(),
	})

	return token.SignedString([]byte(tm.sign))
}

func (tm manager) Parse(token, audience string) (string, error) {
	_, err := url.Parse(audience)
	if err != nil {
		return "", ErrInvalidAudience
	}

	t, err := jwt.Parse(token, func(t *jwt.Token) (i interface{}, err error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSignMethod
		}

		return []byte(tm.sign), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok && !t.Valid {
		return "", ErrNoClaims
	}

	if claims["aud"].(string) != audience {
		return "", ErrAudienceMismatch
	}

	return claims["sub"].(string), nil
}
