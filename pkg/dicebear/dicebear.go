package dicebear

import "fmt"

const (
	avatarSourceTmpl = "https://avatars.dicebear.com/api/identicon/%s.svg"
)

func URL(seed string) string {
	return fmt.Sprintf(avatarSourceTmpl, seed)
}
