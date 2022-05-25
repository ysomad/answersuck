package answer

type (
	CreateRequest struct {
		Answer  string `json:"answer" binding:"required,gte=1,lte=100"`
		MediaId string `json:"mediaId" binding:"uuid4"`
	}
)
