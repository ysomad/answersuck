package packages

type CreateParams struct {
	Name        string
	Description string
	AccountId   string
	LanguageId  uint8
	Tags        []uint32
}
