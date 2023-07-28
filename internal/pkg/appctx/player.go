package appctx

import "context"

type NicknameKey struct{}

// GetNickname returns current user nickname from context.
func GetNickname(ctx context.Context) (string, bool) {
	n, ok := ctx.Value(NicknameKey{}).(string)
	return n, ok
}
