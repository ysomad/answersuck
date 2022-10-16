package dto

type AccountSaveArgs struct {
	Email           string
	Username        string
	EncodedPassword string
}

type AccountCreateArgs struct {
	Email, Username, PlainPassword string
}
