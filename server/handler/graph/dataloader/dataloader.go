package dataloader

import (
	"context"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/usecase"

	"github.com/graph-gophers/dataloader/v7"
	"golang.org/x/xerrors"
)

type Loaders struct {
	AssetCategoryLoader dataloader.Interface[domain.AssetCategoryID, *domain.AssetCategory]
	AssetChangeLoader   dataloader.Interface[domain.RecordID, *AssetChangesAssociation]
	AssetLoader         dataloader.Interface[domain.AssetID, *domain.Asset]
}

func NewLoader(usecase *usecase.Usecase) *Loaders {
	assetCategoryBatcher := &assetCategoryBatcher{usecase: usecase}
	assetChangeBatcher := &assetChangeBatcher{usecase: usecase}
	assetBatcher := &assetBatcher{usecase: usecase}

	return &Loaders{
		AssetCategoryLoader: dataloader.NewBatchedLoader(assetCategoryBatcher.BatchGetAssetCategories),
		AssetChangeLoader:   dataloader.NewBatchedLoader(assetChangeBatcher.BatchGetAssetChanges),
		AssetLoader:         dataloader.NewBatchedLoader(assetBatcher.BatchGetAssets),
	}
}

type assetCategoryBatcher struct {
	usecase *usecase.Usecase
}

func (a *assetCategoryBatcher) BatchGetAssetCategories(ctx context.Context, assetCategoryIDs []domain.AssetCategoryID) []*dataloader.Result[*domain.AssetCategory] {
	results := make([]*dataloader.Result[*domain.AssetCategory], len(assetCategoryIDs))

	indexs := make(map[domain.AssetCategoryID]int, len(assetCategoryIDs))
	for i, ID := range assetCategoryIDs {
		indexs[ID] = i
	}

	userID, err := ctxdef.UserID(ctx)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[*domain.AssetCategory]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	categories, err := a.usecase.GetAssetCategoriesByIDs(ctx, userID, assetCategoryIDs)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[*domain.AssetCategory]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	for _, category := range categories {
		results[indexs[category.ID]] = &dataloader.Result[*domain.AssetCategory]{
			Data:  category,
			Error: nil,
		}
	}

	return results

}
