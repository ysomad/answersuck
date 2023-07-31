package topic

import "time"

type Topic struct {
	ID         int32     `db:"id"`
	Title      string    `db:"title"`
	Author     string    `db:"author"`
	CreateTime time.Time `db:"create_time"`
}
