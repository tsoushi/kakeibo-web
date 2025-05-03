package usecase

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateBank(ctx context.Context, userID domain.UserID, name string) (*domain.Bank, error) {
	bank, err := u.repo.Bank.Insert(ctx, domain.NewBank(userID, name))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return bank, nil
}
