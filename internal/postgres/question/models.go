package question

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype/zeronull"
)

type question struct {
	ID             int32         `db:"id"`
	Text           string        `db:"text"`
	Author         string        `db:"author"`
	MediaURL       zeronull.Text `db:"media_url"`
	CreateTime     time.Time     `db:"create_time"`
	AnswerID       int32         `db:"answer_id"`
	Answer         string        `db:"answer"`
	AnswerMediaURL zeronull.Text `db:"answer_media_url"`
}
