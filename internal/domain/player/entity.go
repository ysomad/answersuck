package player

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound            = errors.New("player not found")
	ErrEmptyAvatarFilename = errors.New("empty avatar filename")
	ErrInvalidAccountId    = errors.New("invalid account id")
)

const (
	MaxAvatarSize = 5 << 20 // 5mb
	dicebearFmt   = "https://avatars.dicebear.com/api/identicon/%s.svg"
)

type Player struct {
	Id             string
	Nickname       string
	Verified       bool
	AvatarFilename *string
}

type Detailed struct {
	Id        string `json:"-"`
	Nickname  string `json:"nickname"`
	Verified  bool   `json:"verified"`
	AvatarURL string `json:"avatar_url"`
}

func NewDetailed(p Player, sp fileStorage) Detailed {
	d := Detailed{
		Id:       p.Id,
		Nickname: p.Nickname,
		Verified: p.Verified,
	}

	if p.AvatarFilename != nil {
		d.AvatarURL = sp.URL(*p.AvatarFilename).String()
	} else {
		d.AvatarURL = fmt.Sprintf(dicebearFmt, d.Nickname)
	}

	return d
}

type Avatar struct {
	AccountId string
	Filename  string
	UpdatedAt time.Time
}
