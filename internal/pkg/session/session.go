package session

import (
	"time"
)

const Cookie = "sid"

type Player struct {
	Nickname   string
	UserAgent  string
	RemoteAddr string
	Verified   bool
}

type Session struct {
	ID        string
	Player    Player
	ExpiresAt time.Time
}

func (s *Session) Expired() bool {
	return s.ExpiresAt.Unix() > time.Now().Unix()
}
