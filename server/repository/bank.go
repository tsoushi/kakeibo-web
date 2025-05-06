package repository

import (
	"context"
	"errors"
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

func (r *BankRepository) List(ctx context.Context, userID domain.UserID) ([]*domain.Bank, error) {
	runner := getRunner(ctx, r.sess)
	banks := make([]*domain.Bank, 0)
	_, err := runner.Select("*").From(banktableName).Where("user_id = ?", userID).LoadContext(ctx, &banks)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to list banks by userID: %w", err)
	}

	return banks, nil
}
