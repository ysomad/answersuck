package account

import "errors"

var (
	ErrAlreadyExist              = errors.New("account with given email or username already exist")
	ErrNotDeleted                = errors.New("account has not been deleted")
	ErrAlreadyArchived           = errors.New("account already archived or not found")
	ErrIncorrectCredentials      = errors.New("incorrect login or password")
	ErrAlreadyVerified           = errors.New("current email already verified or verification code is expired")
	ErrForbiddenNickname         = errors.New("nickname contains forbidden words")
	ErrNotFound                  = errors.New("account not found")
	ErrEmptyVerificationCode     = errors.New("empty account verification code")
	ErrEmptyPasswordResetToken   = errors.New("empty password reset token")
	ErrPasswordTokenNotFound     = errors.New("account password reset token not found or expired")
	ErrIncorrectPassword         = errors.New("incorrect password")
	ErrPasswordResetTokenExpired = errors.New("password reset token is expired")
	ErrVerificationNotFound      = errors.New("account verification not found")
	ErrPasswordNotSet            = errors.New("account password is not set")
	ErrPasswordTokenAlreadyExist = errors.New("account password reset token already exist")
	ErrNotEnoughRights           = errors.New("not enough rights to perform the operation, verify account first")
)
