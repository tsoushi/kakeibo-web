package usecase

import (
	"context"
	"kakeibo-web-server/domain"
)

func (u *Usecase) CreateAssetCategory(ctx context.Context, userID domain.UserID, name string) (*domain.AssetCategory, error) {
	assetCategory, err := u.repo.AssetCategory.Insert(ctx, domain.NewAssetCategory(userID, name))
	if err != nil {
		return nil, err
	}

	return assetCategory, nil
}

func (u *Usecase) GetAssetCategoryByID(ctx context.Context, userID domain.UserID, id domain.AssetCategoryID) (*domain.AssetCategory, error) {
	category, err := u.repo.AssetCategory.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (u *Usecase) GetAssetCategoriesByUserID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID) ([]*domain.AssetCategory, *domain.PageInfo, error) {
	categories, pageInfo, err := u.repo.AssetCategory.GetMultiByUserID(ctx, pageParam, userID)
	if err != nil {
		return nil, nil, err
	}

	return categories, pageInfo, nil
}
