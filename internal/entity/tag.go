package entity

import (
	"time"
)

const TagPageSize = 2

type Tag struct {
	Name       string
	Author     string
	CreateTime time.Time
}
