package account

import (
	"time"
)

type (
	CreateReq struct {
		Email    string `json:"email" validate:"required,email,lte=255"`
		Nickname string `json:"nickname" validate:"required,alphanum,gte=4,lte=25"`
		Password string `json:"password" validate:"required,gte=8,lte=71"`
	}

	ResetPasswordReq struct {
		Login string `json:"login" binding:"required,email|alphanum"`
	}

	SetPasswordReq struct {
		Password string `json:"password" binding:"required,gte=8,lte=71"`
	}
)

type (
	VerificationDTO struct {
		Email    string
		Code     string
		Verified bool
	}

	SetPasswordDTO struct {
		AccountId string
		Password  string
		Token     string
		UpdatedAt time.Time
	}
)
