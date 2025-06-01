package usecase

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateTag(ctx context.Context, userID domain.UserID, name string) (*domain.Tag, error) {
	createdTag, err := u.repo.Tag.Insert(ctx, domain.NewTag(userID, name))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return createdTag, nil
}
