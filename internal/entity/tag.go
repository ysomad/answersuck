package entity

import (
	"time"

	"github.com/ysomad/answersuck/internal/pkg/paging"
)

var _ paging.SeekListItem = &Tag{}

type Tag struct {
	Name       string
	Author     string
	CreateTime time.Time
}

func (t Tag) GetID() string      { return t.Name }
func (t Tag) GetTime() time.Time { return t.CreateTime }
