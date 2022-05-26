package answer

type (
	CreateRequest struct {
		Text    string `json:"text" binding:"required,gte=1,lte=100"`
		MediaId string `json:"mediaId" binding:"omitempty,uuid4"`
	}
)
