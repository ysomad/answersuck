package dto

type (
	LoginRequest struct {
		Login    string `json:"login" binding:"required,email|alphanum"`
		Password string `json:"password" binding:"required"`
	}

	TokenCreateRequest struct {
		Audience string `json:"audience" binding:"required,uri"`
		Password string `json:"password" binding:"required"`
	}

	TokenCreateResponse struct {
		Token string `json:"token"`
	}
)
