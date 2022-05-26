package session

type (
	SessionWithVerified struct {
		Session         Session
		AccountVerified bool
	}
)
