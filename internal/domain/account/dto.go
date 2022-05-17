package account

import (
	"time"
)

type (
	CreateRequest struct {
		Email    string `json:"email" binding:"required,email,lte=255"`
		Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
		Password string `json:"password" binding:"required,gte=8,lte=64"`
	}

	ResetPasswordRequest struct {
		Login string `json:"login" binding:"required,email|alphanum"`
	}

	SetPasswordRequest struct {
		Password string `json:"password" binding:"required,gte=8,lte=64"`
	}
)

type (
	VerificationDTO struct {
		Email    string
		Code     string
		Verified bool
	}

	SetPasswordDTO struct {
		AccountId    string
		Token        string
		PasswordHash string
		UpdatedAt    time.Time
	}
)
