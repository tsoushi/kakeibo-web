package usecase

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateAssetCategory(ctx context.Context, userID domain.UserID, name string) (*domain.AssetCategory, error) {
	assetCategory, err := u.repo.AssetCategory.Insert(ctx, domain.NewAssetCategory(userID, name))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return assetCategory, nil
}

func (u *Usecase) GetAssetCategoryByID(ctx context.Context, userID domain.UserID, id domain.AssetCategoryID) (*domain.AssetCategory, error) {
	category, err := u.repo.AssetCategory.GetByID(ctx, userID, id)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return category, nil
}

func (u *Usecase) GetAssetCategoriesByUserID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID) ([]*domain.AssetCategory, *domain.PageInfo, error) {
	categories, pageInfo, err := u.repo.AssetCategory.GetMultiByUserID(ctx, pageParam, userID)
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return categories, pageInfo, nil
}

func (u *Usecase) GetAssetCategoriesByIDs(ctx context.Context, userID domain.UserID, assetCategoryIDs []domain.AssetCategoryID) ([]*domain.AssetCategory, error) {
	categories, err := u.repo.AssetCategory.GetMultiByAssetCategoryIDs(ctx, userID, assetCategoryIDs)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return categories, nil
}
