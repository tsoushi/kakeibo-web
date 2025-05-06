package domain

import "time"

const (
	BankIDSuffix = "Bank"
)

type BankID string

func NewBankID() BankID {
	return BankID(NewUUIDv4(BankIDSuffix))
}

type Bank struct {
	ID        BankID
	UserID    UserID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBank(userID UserID, name string) *Bank {
	return &Bank{
		ID:        NewBankID(),
		UserID:    userID,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
