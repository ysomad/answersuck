package session

import (
	"net"
	"time"
)

const Cookie = "sid"

type User struct {
	UserID    string
	UserAgent string
	IP        net.IP
	Verified  bool
}

type Session struct {
	ID        string
	User      User
	ExpiresAt time.Time
}

func (s *Session) Expired() bool {
	return s.ExpiresAt.Unix() > time.Now().Unix()
}
