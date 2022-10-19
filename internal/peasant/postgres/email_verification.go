package postgres

import (
	"context"

	"github.com/ysomad/answersuck/apperror"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/pgclient"
)

type emailVerificationRepository struct {
	*pgclient.Client
}

func NewEmailVerificationRepository(c *pgclient.Client) *emailVerificationRepository {
	return &emailVerificationRepository{c}
}

func (s *emailVerificationRepository) Save(ctx context.Context, v domain.EmailVerification) error {
	query, args, err := s.Builder.
		Insert("email_verification").
		Columns("account_id, code, expires_at").
		Values(v.AccountID, v.Code, v.ExpiresAt).
		ToSql()
	if err != nil {
		return nil
	}

	ct, err := s.Pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return apperror.New("emailVerificationRepository - Save", apperror.ErrZeroRowsAffected, domain.ErrEmailVerificationNotCreated)
	}

	return nil
}
