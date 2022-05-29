package session

type (
	WithAccountDetails struct {
		Session  Session
		Verified bool
	}
)
