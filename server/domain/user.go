package domain

import (
	"time"
)

const (
	UserIDSuffix = "User"
)

type UserID string

func NewUserID() UserID {
	return UserID(NewUUIDv4(UserIDSuffix))
}

type User struct {
	ID             UserID
	Name           string
	HashedPassword HashedPassword
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUser(name string, hashedPassword HashedPassword) *User {
	return &User{
		ID:             NewUserID(),
		Name:           name,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
