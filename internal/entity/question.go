package entity

import "time"

type Answer struct {
	ID       int32
	Text     string
	MediaURL string
}

type Question struct {
	ID         int32
	Text       string
	Answer     Answer
	Author     string
	MediaURL   string
	CreateTime time.Time
}
