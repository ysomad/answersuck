package entity

import "time"

type Account struct {
	ID            string
	Email         string
	Username      string
	PasswordHash  string
	EmailVerified bool
	Archived      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (a *Account) BeforeSave(loc *time.Location) {
	now := time.Now().In(loc)
	a.CreatedAt = now
	a.UpdatedAt = now
}

func (a *Account) BeforeUpdate(loc *time.Location) {
	a.UpdatedAt = time.Now().In(loc)
}
