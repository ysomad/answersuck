package session

import (
	"net"
	"time"
)

const Cookie = "sid"

type Player struct {
	Nickname  string
	UserAgent string
	IP        net.IP
	Verified  bool
}

type Session struct {
	ID        string
	Player    Player
	ExpiresAt time.Time
}

func (s *Session) Expired() bool {
	return s.ExpiresAt.Unix() > time.Now().Unix()
}
