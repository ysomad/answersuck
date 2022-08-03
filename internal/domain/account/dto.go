package account

import (
	"time"
)

type (
	CreateReq struct {
		Email    string `json:"email" validate:"required,email,lte=255"`
		Nickname string `json:"nickname" validate:"required,alphanum,gte=4,lte=25"`
		Password string `json:"password" validate:"required,gte=10,lte=128"`
	}

	UpdatePasswordReq struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,gte=10,lte=128,nefield=OldPassword"`
	}

	ResetPasswordReq struct {
		Login string `json:"login" validate:"required,email|alphanum"`
	}

	SetPasswordReq struct {
		Token    string `json:"token" validate:"required"`
		Password string `json:"password" validate:"required,gte=10,lte=128"`
	}
)

type (
	SetPasswordDTO struct {
		AccountId string
		Password  string
		Token     string
		UpdatedAt time.Time
	}

	SavePasswordTokenDTO struct {
		Login     string
		Token     string
		CreatedAt time.Time
	}
)
