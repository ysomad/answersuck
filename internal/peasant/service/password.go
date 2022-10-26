package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	"github.com/ysomad/answersuck/jwt"
)

// passwordService encapsulates logic related to account password.
type passwordService struct {
	accountRepo accountRepository
	password    passwordEncodeComparer

	setterToken    basicJWTManager
	setterTokenExp time.Duration
}

func NewPasswordService(r accountRepository, ec passwordEncodeComparer, m basicJWTManager, setterTokenExp time.Duration) *passwordService {
	return &passwordService{
		accountRepo:    r,
		password:       ec,
		setterToken:    m,
		setterTokenExp: setterTokenExp,
	}
}

// Update updates account password if old password is correct.
func (s *passwordService) Update(ctx context.Context, args dto.UpdatePasswordArgs) (*domain.Account, error) {
	oldEncodedPass, err := s.accountRepo.GetPasswordByID(ctx, args.AccountID)
	if err != nil {
		return nil, err
	}

	ok, err := s.password.Compare(args.OldPassword, oldEncodedPass)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, domain.ErrIncorrectPassword
	}

	newEncodedPass, err := s.password.Encode(args.NewPassword)
	if err != nil {
		return nil, err
	}

	return s.accountRepo.UpdatePassword(ctx, args.AccountID, newEncodedPass)
}

// NotifyWithToken finds account with email or username and creates password setter token with found account id,
// then send email to user with token in url, user must to visit the url to set new password.
func (s *passwordService) NotifyWithToken(ctx context.Context, emailOrUsername string) (domain.PasswordSetterToken, error) {
	// TODO: get account by email or username
	// TODO: implement password service create token

	return "", nil
}

// Set sets new password to account. Token must be valid not expired jwt with account id in payload.
func (s *passwordService) Set(ctx context.Context, token, newPassword string) (*domain.Account, error) {
	c, err := s.setterToken.Decode(jwt.Basic(token))
	if err != nil {
		return nil, apperror.New("passwordService - Set", err, domain.ErrPasswordSetterTokenExpired)
	}

	newEncodedPass, err := s.password.Encode(newPassword)
	if err != nil {
		return nil, err
	}

	return s.accountRepo.UpdatePassword(ctx, c.Subject, newEncodedPass)
}
