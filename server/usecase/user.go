package usecase

import (
	"context"
	"errors"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) GetUserByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	user, err := u.repo.User.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrEntityNotFound) {
			return nil, xerrors.Errorf(": %w", err)
		}
		return nil, xerrors.Errorf(": %w", err)
	}

	return user, nil
}
