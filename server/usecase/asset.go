package usecase

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateAsset(ctx context.Context, userID domain.UserID, name string, categoryID *domain.AssetCategoryID) (*domain.Asset, error) {
	asset, err := u.repo.Asset.Insert(ctx, domain.NewAsset(userID, name, categoryID))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return asset, nil
}

func (u *Usecase) UpdateAsset(ctx context.Context, userID domain.UserID, id domain.AssetID, name string, categoryID *domain.AssetCategoryID) (*domain.Asset, error) {
	asset := &domain.Asset{
		ID:         id,
		UserID:     userID,
		Name:       name,
		CategoryID: categoryID,
	}

	updatedAsset, err := u.repo.Asset.Update(ctx, asset)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return updatedAsset, nil
}

func (u *Usecase) DeleteAsset(ctx context.Context, userID domain.UserID, id domain.AssetID) (domain.AssetID, error) {
	asset := &domain.Asset{
		ID:     id,
		UserID: userID,
	}

	deletedAsset, err := u.repo.Asset.Delete(ctx, asset)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}

	return deletedAsset, nil
}

func (u *Usecase) GetAssetsByUserIDAndCategoryID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, categoryID *domain.AssetCategoryID) ([]*domain.Asset, *domain.PageInfo, error) {
	assets, pageInfo, err := u.repo.Asset.GetMultiByUserIDAndCategoryID(ctx, pageParam, userID, categoryID)
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return assets, pageInfo, nil
}

func (u *Usecase) GetAssetsByIDs(ctx context.Context, userID domain.UserID, assetIDs []domain.AssetID) ([]*domain.Asset, error) {
	if len(assetIDs) == 0 {
		return nil, nil
	}

	assets, err := u.repo.Asset.GetMultiByUserIDAndIDs(ctx, userID, assetIDs)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return assets, nil
}

func (u *Usecase) GetAssetsByCategoryIDs(ctx context.Context, userID domain.UserID, assetCategoryIDs []domain.AssetCategoryID) ([]*domain.Asset, error) {
	if len(assetCategoryIDs) == 0 {
		return nil, nil
	}

	assets, err := u.repo.Asset.GetMultiByCategoryIDs(ctx, userID, assetCategoryIDs)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return assets, nil
}
