package domain

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

const (
	BCRYPT_COST = 10
)

type HashedPassword string

func NewHashedPassword(plainPassword string) (HashedPassword, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(plainPassword), BCRYPT_COST)
	if err != nil {
		return "", err
	}
	return HashedPassword(hashedPassword), nil
}

func (h HashedPassword) Compare(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
