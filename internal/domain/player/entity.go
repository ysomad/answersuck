package player

import "errors"

var (
	ErrNotFound = errors.New("player not found")
)

type Player struct {
	Id        string  `json:"id"`
	Nickname  string  `json:"nickname"`
	AvatarURL *string `json:"avatar_url"`
}

func (p *Player) setAvatarURL(mp mediaProvider) {
	if p.AvatarURL != nil {
		avatarURL := mp.URL(*p.AvatarURL).String()
		p.AvatarURL = &avatarURL
	}
}
