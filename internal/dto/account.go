package dto

type AccountVerification struct {
	Email    string
	Code     string
	Verified bool
}
