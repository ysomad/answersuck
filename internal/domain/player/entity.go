package player

type Player struct {
	Id        string  `json:"-"`
	Nickname  string  `json:"nickname"`
	AvatarURL *string `json:"avatarUrl"`
}
