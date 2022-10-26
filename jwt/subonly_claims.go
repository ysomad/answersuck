package jwt

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type SubOnlyClaims struct {
	ExpiresAt int64  `json:"exp"`
	IssuedAt  int64  `json:"iat"`
	Subject   string `json:"sub"`
	Issuer    string `json:"iss"`
}

func NewSubOnlyClaims(subject, issuer string, exp time.Duration) (SubOnlyClaims, error) {
	if exp <= 0 {
		return SubOnlyClaims{}, errInvalidExpiration
	}

	if issuer == "" {
		return SubOnlyClaims{}, errInvalidIssuer
	}

	_, err := uuid.Parse(subject)
	if err != nil {
		return SubOnlyClaims{}, fmt.Errorf("%s: %w", err.Error(), errInvalidSubject)
	}

	now := time.Now()
	return SubOnlyClaims{
		ExpiresAt: now.Add(exp).Unix(),
		IssuedAt:  now.Unix(),
		Subject:   subject,
		Issuer:    issuer,
	}, nil
}

func newSubOnlyClaims(raw jwt.MapClaims) (SubOnlyClaims, error) {
	b, err := json.Marshal(raw)
	if err != nil {
		return SubOnlyClaims{}, err
	}

	var c SubOnlyClaims
	if err := json.Unmarshal(b, &c); err != nil {
		return SubOnlyClaims{}, err
	}

	return c, nil
}

func (c SubOnlyClaims) Valid() error {
	vErr := new(jwt.ValidationError)
	now := time.Now().Unix()

	if !c.verifyExpiresAt(now) {
		delta := time.Unix(now, 0).Sub(time.Unix(c.ExpiresAt, 0))
		vErr.Inner = fmt.Errorf("token is expired by %v", delta)
		vErr.Errors |= jwt.ValidationErrorExpired
	}

	if !c.verifySubject() {
		vErr.Inner = errInvalidSubject
		vErr.Errors |= jwt.ValidationErrorClaimsInvalid
	}

	if !c.verifyIssuer() {
		vErr.Inner = errInvalidIssuer
		vErr.Errors |= jwt.ValidationErrorIssuer
	}

	if !c.verifyIssuedAt(now) {
		vErr.Inner = errors.New("token used before issued")
		vErr.Errors |= jwt.ValidationErrorIssuedAt
	}

	if vErr.Errors == 0 {
		return nil
	}

	return vErr
}

func (c SubOnlyClaims) verifyExpiresAt(now int64) bool {
	if c.ExpiresAt == 0 {
		return false
	}
	return now <= c.ExpiresAt
}

func (c SubOnlyClaims) verifySubject() bool {
	if c.Subject == "" {
		return false
	}

	_, err := uuid.Parse(c.Subject)
	return err == nil
}

func (c SubOnlyClaims) verifyIssuer() bool { return c.Issuer != "" }

func (c SubOnlyClaims) verifyIssuedAt(now int64) bool {
	if c.IssuedAt == 0 {
		return false
	}
	return now >= c.IssuedAt
}
