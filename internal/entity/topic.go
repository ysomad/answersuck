package entity

import "time"

type Topic struct {
	ID         int32
	Title      string
	Author     string
	CreateTime time.Time
}
