package account

import (
	"time"
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
