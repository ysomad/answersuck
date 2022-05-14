package dto

import (
	"time"
)

type (
	AccountCreateRequest struct {
		Email    string `json:"email" binding:"required,email,lte=255"`
		Username string `json:"username" binding:"required,alphanum,gte=4,lte=16"`
		Password string `json:"password" binding:"required,gte=8,lte=64"`
	}

	AccountCreateResponse struct {
		Id        string    `json:"id"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		Verified  bool      `json:"verified"`
		Archived  bool      `json:"archived"`
		AvatarURL string    `json:"avatarUrl"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	AccountPasswordResetRequest struct {
		Login string `json:"login" binding:"required,email|alphanum"`
	}

	AccountPasswordSetRequest struct {
		Password string `json:"password" binding:"required,gte=8,lte=64"`
	}
)

type (
	AccountVerification struct {
		Email    string
		Code     string
		Verified bool
	}

	AccountPasswordResetToken struct {
		AccountId string
		Token     string
		CreatedAt time.Time
	}

	AccountUpdatePassword struct {
		AccountId    string
		Token        string
		PasswordHash string
		UpdatedAt    time.Time
	}
)
