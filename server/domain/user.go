package domain

import "time"

type UserID string

type User struct {
	ID        UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
