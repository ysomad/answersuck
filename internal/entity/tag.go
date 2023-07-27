package entity

import (
	"time"
)

type Tag struct {
	Name       string
	Author     string
	CreateTime time.Time
}
