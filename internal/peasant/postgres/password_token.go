package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ysomad/answersuck/internal/peasant/domain"
	"github.com/ysomad/answersuck/internal/peasant/service/dto"
	"github.com/ysomad/answersuck/pgclient"
)

type passwordTokenRepository struct {
	*pgclient.Client
}

func NewPasswordTokenRepository(c *pgclient.Client) *passwordTokenRepository {
	return &passwordTokenRepository{c}
}

func (r *passwordTokenRepository) Create(ctx context.Context, args dto.CreatePasswordTokenArgs) (domain.PasswordToken, error) {
	query := `
INSERT INTO password_token(account_id, token, expires_at)
VALUES ((SELECT id FROM account WHERE email = $1 OR username = $2), $3, $4)
RETURNING account_id, token, expires_at`

	rows, err := r.Pool.Query(
		ctx,
		query,
		args.EmailOrUsername,
		args.EmailOrUsername,
		args.Token,
		args.ExpiresAt,
	)
	if err != nil {
		return domain.PasswordToken{}, err
	}

	t, err := pgx.CollectOneRow(rows, pgx.RowToStructByPos[domain.PasswordToken])
	if err != nil {
		return domain.PasswordToken{}, err
	}

	return t, nil
}
