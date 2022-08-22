package psql

import (
	"go.uber.org/zap"

	"github.com/ysomad/answersuck-backend/internal/pkg/postgres"
)

type PackageRepo struct {
	*zap.Logger
	*postgres.Client
}

func NewPackageRepo(l *zap.Logger, c *postgres.Client) *PackageRepo {
	return &PackageRepo{l, c}
}

// func (r *PackageRepo) Save(ctx context.Context, p packages.Package) (packageId uint32, err error) {
// 	sql, args, err := r.Builder.
// 		Insert("package").
// 		Columns("name, description, account_id, is_published, language_id, created_at, updated_at").
// 		Values(p.Name, p.Description, p.AccountId, p.Published, p.LanguageId, p.CreatedAt, p.UpdatedAt).
// 		Suffix("RETURNING id").
// 		ToSql()
// 	if err != nil {
// 		return 0, fmt.Errorf("psql - package - Save - ToSql: %w", err)
// 	}
//
//     ib := r.Builder.
//         Insert("package_tag").
//         Columns("package_id, tag_id").
//
// 	sql := `WITH p AS (
// INSERT INTO package(name, description, account_id, is_published, language_id, created_at, updated_at)
// VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
// ), t AS (
//     INSERT INTO package_tag
// )`
//
// 	r.Debug("psql - package - Save", zap.String("sql", sql), zap.Any("args", args))
//
// 	return 0, nil
// }
