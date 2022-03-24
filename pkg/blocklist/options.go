package blocklist

type option func() []string

func WithUsernames() []string {
	return Usernames
}
