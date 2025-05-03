package repository

import (
	"context"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const banktableName = "bank"

type BankRepository struct {
	sess *dbr.Session
}

func NewBankRepository(sess *dbr.Session) *BankRepository {
	return &BankRepository{
		sess: sess,
	}
}

func (r *BankRepository) Insert(ctx context.Context, bank *domain.Bank) (*domain.Bank, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(banktableName).Columns("id", "user_id", "name").Record(bank).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert bank: %w", err)
	}

	return bank, nil
}
