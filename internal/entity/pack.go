package entity

import "time"

type Pack struct {
	ID         int32
	Name       string
	Author     string
	Published  bool
	CoverURL   string
	CreateTime time.Time
}
