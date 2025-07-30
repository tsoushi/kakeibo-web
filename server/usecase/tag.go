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

func (u *Usecase) UpdateTag(ctx context.Context, userID domain.UserID, id domain.TagID, name string) (*domain.Tag, error) {
	tag := &domain.Tag{
		ID:     id,
		UserID: userID,
		Name:   name,
	}
	updatedTag, err := u.repo.Tag.Update(ctx, tag)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return updatedTag, nil
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

func (u *Usecase) GetOrCreateTagsByName(ctx context.Context, userID domain.UserID, names []string) ([]*domain.Tag, error) {
	tags := make([]*domain.Tag, 0, len(names))

	existTags, err := u.repo.Tag.GetMultiByNames(ctx, userID, names)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	tags = append(tags, existTags...)

	newTags := domain.NewTagsNotExist(userID, existTags, names)
	for _, tag := range newTags {
		createdTag, err := u.repo.Tag.Insert(ctx, tag)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		tags = append(tags, createdTag)
	}

	return tags, nil
}

func (u *Usecase) GetTagsWithRecordIDByRecordIDs(ctx context.Context, userID domain.UserID, recordIDs []domain.RecordID) ([]*domain.TagWithRecordID, error) {
	if len(recordIDs) == 0 {
		return nil, nil
	}

	tagWithRecordIDs, err := u.repo.Tag.GetMultiWithRecordIDByRecordIDs(ctx, userID, recordIDs)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return tagWithRecordIDs, nil
}
