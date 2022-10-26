package service

import "github.com/ysomad/answersuck/jwt"

type passwordEncoder interface {
	Encode(plain string) (string, error)
}

type passwordComparer interface {
	Compare(plain, encoded string) (bool, error)
}

type passwordEncodeComparer interface {
	passwordEncoder
	passwordComparer
}

type basicJWTManager interface {
	Encode(jwt.BasicClaims) (jwt.Basic, error)
	Decode(jwt.Basic) (jwt.BasicClaims, error)
	Issuer() string
}
