package pack

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
)

type Pack struct {
	ID         int32         `db:"id"`
	Name       string        `db:"name"`
	Author     string        `db:"author"`
	Published  bool          `db:"is_published"`
	CoverURL   zeronull.Text `db:"cover_url"`
	CreateTime time.Time     `db:"create_time"`
}
