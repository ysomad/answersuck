package auth

type (
	LoginReq struct {
		Login    string `json:"login" validate:"required,email|alphanum"`
		Password string `json:"password" validate:"required"`
	}

	TokenCreateReq struct {
		Password string `json:"password" validate:"required"`
	}

	TokenCreateResp struct {
		Token string `json:"token"`
	}
)
