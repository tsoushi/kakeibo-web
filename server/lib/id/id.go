package id

import "github.com/google/uuid"

func NewUUIDv4() string {
	return uuid.New().String()
}
