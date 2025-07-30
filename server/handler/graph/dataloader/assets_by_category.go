package dataloader

import (
	"context"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/usecase"

	"github.com/graph-gophers/dataloader/v7"
	"golang.org/x/xerrors"
)

type assetsByCategoryBatcher struct {
	usecase *usecase.Usecase
}

func (a *assetsByCategoryBatcher) BatchGetAssetsByCategoryIDs(ctx context.Context, assetCategoryIDs []domain.AssetCategoryID) []*dataloader.Result[[]*domain.Asset] {
	results := make([]*dataloader.Result[[]*domain.Asset], len(assetCategoryIDs))

	indexs := make(map[domain.AssetCategoryID]int, len(assetCategoryIDs))
	for i, ID := range assetCategoryIDs {
		indexs[ID] = i
	}

	userID, err := ctxdef.UserID(ctx)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[[]*domain.Asset]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	assets, err := a.usecase.GetAssetsByCategoryIDs(ctx, userID, assetCategoryIDs)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[[]*domain.Asset]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	for _, assetCategoryID := range assetCategoryIDs {
		results[indexs[assetCategoryID]] = &dataloader.Result[[]*domain.Asset]{
			Data:  make([]*domain.Asset, 0),
			Error: nil,
		}
	}

	for _, asset := range assets {
		if asset.CategoryID == nil {
			panic("asset.CategoryID is nil")
		}
		results[indexs[*asset.CategoryID]].Data = append(results[indexs[*asset.CategoryID]].Data, asset)
	}

	return results
}
