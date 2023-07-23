package player

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
)

type player struct {
	Nickname      string               `db:"nickname"`
	Email         string               `db:"email"`
	DisplayName   zeronull.Text        `db:"display_name"`
	EmailVerified bool                 `db:"email_verified"`
	Password      string               `db:"password"`
	CreateTime    time.Time            `db:"create_time"`
	UpdateTime    zeronull.Timestamptz `db:"update_time"`
}
