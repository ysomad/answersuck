package domain

type Tag struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	LanguageId int    `json:"languageId"`
}
