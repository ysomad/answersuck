package dto

type AccountSaveArgs struct {
	Email           string
	Username        string
	EncodedPassword string
	EmailVerifCode  string
}
