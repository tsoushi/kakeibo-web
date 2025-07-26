package domain

import (
	"time"
)

const (
	UserIDSuffix = "User"
)

type UserID string

type User struct {
	ID        UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(id UserID, name string) *User {
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
