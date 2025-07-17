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

func (u *Usecase) GetTagsByUserID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID) ([]*domain.Tag, *domain.PageInfo, error) {
	tags, pageInfo, err := u.repo.Tag.GetMultiByUserID(ctx, pageParam, userID)
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return tags, pageInfo, nil
}

func (u *Usecase) DeleteTag(ctx context.Context, userID domain.UserID, tagID domain.TagID) (domain.TagID, error) {
	deletedTag, err := u.repo.Tag.Delete(ctx, userID, tagID)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}

	return deletedTag, nil
}
