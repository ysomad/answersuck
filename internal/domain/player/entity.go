package player

type Player struct {
	Id        string  `json:"id"`
	Nickname  string  `json:"nickname"`
	AvatarURL *string `json:"avatarUrl"`
}
