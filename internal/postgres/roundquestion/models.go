package roundquestion

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgtype/zeronull"

	"github.com/ysomad/answersuck/internal/entity"
)

type roundQuestion struct {
	ID           int32               `db:"id"`
	TopicID      int32               `db:"topic_id"`
	RoundID      int32               `db:"round_id"`
	Type         entity.QuestionType `db:"question_type"`
	Cost         int32               `db:"question_cost"`
	AnswerTime   time.Duration       `db:"answer_time"`
	HostComment  zeronull.Text       `db:"host_comment"`
	SecretTopic  zeronull.Text       `db:"secret_topic"`
	SecretCost   zeronull.Int2       `db:"secret_cost"`
	Keepable     pgtype.Bool         `db:"is_keepable"`
	TransferType zeronull.Int2       `db:"transfer_type"`

	QuestionID       int32         `db:"question_id"`
	Question         string        `db:"question"`
	QuestionMediaURL zeronull.Text `db:"question_media_url"`

	AnswerID       int32         `db:"answer_id"`
	Answer         string        `db:"answer"`
	AnswerMediaURL zeronull.Text `db:"answer_media_url"`
}
