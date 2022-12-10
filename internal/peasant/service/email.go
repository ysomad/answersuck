package service

import (
	"context"
	"time"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	"github.com/ysomad/answersuck/jwt"
)

type emailService struct {
	accountRepo accountRepository
	password    passwordComparer

	verifToken    basicJWTManager
	verifTokenExp time.Duration
}

func NewEmailService(r accountRepository, p passwordComparer, m basicJWTManager, verifTokenExp time.Duration) *emailService {
	return &emailService{
		accountRepo:   r,
		password:      p,
		verifToken:    m,
		verifTokenExp: verifTokenExp,
	}
}

func (s *emailService) Update(ctx context.Context, args dto.UpdateEmailArgs) (*domain.Account, error) {
	encoded, err := s.accountRepo.GetPasswordByID(ctx, args.AccountID)
	if err != nil {
		return nil, err
	}

	ok, err := s.password.Compare(args.PlainPassword, encoded)
	if err != nil {
		return nil, apperror.New("emailService - Update", err, domain.ErrIncorrectPassword)
	}
	if !ok {
		return nil, domain.ErrIncorrectPassword
	}

	a, err := s.accountRepo.UpdateEmail(ctx, args.AccountID, args.NewEmail)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Verify verifies account email if token is valid and contains account id,
// account with already verified email cannot be verified.
func (s *emailService) Verify(ctx context.Context, token string) (*domain.Account, error) {
	c, err := s.verifToken.Decode(token)
	if err != nil {
		return nil, apperror.New("emailService - Verify", err, domain.ErrEmailVerifTokenExpired)
	}

	return s.accountRepo.VerifyEmail(ctx, c.Subject)
}

// NotifyWithToken creates new email verification token and notifies user with it in url, user must visit
// the url to verify email.
func (s *emailService) NotifyWithToken(ctx context.Context, accountID string) (domain.EmailVerifToken, error) {
	a, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		return "", err
	}

	c := jwt.NewBasicClaims(a.ID, s.verifToken.Issuer(), s.verifTokenExp)
	t, err := s.verifToken.Encode(c)
	if err != nil {
		return "", err
	}

	// TODO: send email verif token to user email
	// emailInterctor.Send(a.Email)

	return domain.EmailVerifToken(t), nil
}
