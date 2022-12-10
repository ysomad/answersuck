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
	Encode(jwt.BasicClaims) (string, error)
	Decode(string) (jwt.BasicClaims, error)
	Issuer() string
}
