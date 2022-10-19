package dto

type UpdatePasswordArgs struct {
	AccountID   string
	OldPassword string
	NewPassword string
}
