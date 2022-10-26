package jwt

import "errors"

var (
	errEmptySign            = errors.New("signing key is empty")
	errNoClaims             = errors.New("error getting claims from token")
	errUnexpectedSignMethod = errors.New("unexpected signing method")
	errInvalidExpiration    = errors.New("expiration must be greater than 0")
	errInvalidIssuer        = errors.New("issuer cannot be empty string")
	errInvalidSubject       = errors.New("subject must be valid uuid")
)
