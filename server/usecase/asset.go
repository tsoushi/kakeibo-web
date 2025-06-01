package usecase

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateAsset(ctx context.Context, userID domain.UserID, name string) (*domain.Asset, error) {
	asset, err := u.repo.Asset.Insert(ctx, domain.NewAsset(userID, name))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return asset, nil
}

func (u *Usecase) GetAssetsByUserID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID) ([]*domain.Asset, *domain.PageInfo, error) {
	assets, pageInfo, err := u.repo.Asset.GetMultiByUserID(ctx, pageParam, userID)
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return assets, pageInfo, nil
}
