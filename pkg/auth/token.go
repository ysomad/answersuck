package auth

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

type TokenManager interface {
	// New creates new jwt token with subject and audience claims
	New(subject, audience string, ttl time.Duration) (string, error)

	// Parse parses jwt token and returns subject id from payload
	Parse(token, audience string) (string, error)
}

type tokenManager struct {
	sign string
}

func NewTokenManager(sign string) (tokenManager, error) {
	if sign == "" {
		return tokenManager{}, ErrEmptySign
	}

	return tokenManager{
		sign: sign,
	}, nil
}

func (tm tokenManager) New(subject, audience string, ttl time.Duration) (string, error) {
	_, err := url.Parse(audience)
	if err != nil {
		return "", ErrInvalidAudience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   subject,
		Audience:  audience,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	})

	return token.SignedString([]byte(tm.sign))
}

func (tm tokenManager) Parse(token, audience string) (string, error) {
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
