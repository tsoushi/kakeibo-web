package dataloader

import (
	"context"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/usecase"

	"github.com/graph-gophers/dataloader/v7"
	"golang.org/x/xerrors"
)

type assetChangeBatcher struct {
	usecase *usecase.Usecase
}

type AssetChangesAssociation struct {
	RecordID           domain.RecordID
	AssetChangeExpense *domain.AssetChange
	AssetChangeIncome  *domain.AssetChange
}

func (a *assetChangeBatcher) BatchGetAssetChanges(ctx context.Context, recordIDs []domain.RecordID) []*dataloader.Result[*AssetChangesAssociation] {
	results := make([]*dataloader.Result[*AssetChangesAssociation], len(recordIDs))

	indexs := make(map[domain.RecordID]int, len(recordIDs))
	for i, ID := range recordIDs {
		indexs[ID] = i
	}

	userID, err := ctxdef.UserID(ctx)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[*AssetChangesAssociation]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	assetChanges, err := a.usecase.GetAssetChangesByRecordIDs(ctx, userID, recordIDs)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[*AssetChangesAssociation]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	for _, recordID := range recordIDs {
		results[indexs[recordID]] = &dataloader.Result[*AssetChangesAssociation]{
			Data: &AssetChangesAssociation{
				RecordID:           recordID,
				AssetChangeExpense: nil,
				AssetChangeIncome:  nil,
			},
			Error: nil,
		}
	}

	for _, change := range assetChanges {
		if change.Amount >= 0 {
			results[indexs[change.RecordID]].Data.AssetChangeIncome = change
		} else {
			results[indexs[change.RecordID]].Data.AssetChangeExpense = change
		}
	}

	return results
}
