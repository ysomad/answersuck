package session

type (
	WithAccountDetails struct {
		Session  Session
		Nickname string
		Verified bool
	}
)
