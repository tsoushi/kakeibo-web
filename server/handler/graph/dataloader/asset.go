package dataloader

import (
	"context"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/usecase"

	"github.com/graph-gophers/dataloader/v7"
	"golang.org/x/xerrors"
)

type assetBatcher struct {
	usecase *usecase.Usecase
}

func (a *assetBatcher) BatchGetAssets(ctx context.Context, assetIDs []domain.AssetID) []*dataloader.Result[*domain.Asset] {
	results := make([]*dataloader.Result[*domain.Asset], len(assetIDs))

	indexs := make(map[domain.AssetID]int, len(assetIDs))
	for i, ID := range assetIDs {
		indexs[ID] = i
	}

	userID, err := ctxdef.UserID(ctx)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[*domain.Asset]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	assets, err := a.usecase.GetAssetsByIDs(ctx, userID, assetIDs)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[*domain.Asset]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	for _, asset := range assets {
		results[indexs[asset.ID]] = &dataloader.Result[*domain.Asset]{
			Data:  asset,
			Error: nil,
		}
	}

	return results
}
