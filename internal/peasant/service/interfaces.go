package service

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
