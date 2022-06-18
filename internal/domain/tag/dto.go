package tag

type (
	CreateRequest struct {
		Name       string `json:"name" binding:"required,gte=1,lte=32"`
		LanguageId int    `json:"languageId" binding:"required"`
	}

	CreateMultipleReq struct {
		Tags []CreateRequest `json:"tags" binding:"required,min=1,max=10"`
	}
)
