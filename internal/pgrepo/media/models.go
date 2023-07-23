package media

import (
	"time"

	"github.com/ysomad/answersuck/internal/entity"
)

type media struct {
	URL        string           `db:"url"`
	Type       entity.MediaType `db:"type"`
	Author     string           `db:"author"`
	CreateTime time.Time        `db:"create_time"`
}
