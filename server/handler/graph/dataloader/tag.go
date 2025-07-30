package dataloader

import (
	"context"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/usecase"

	"github.com/graph-gophers/dataloader/v7"
	"golang.org/x/xerrors"
)

type tagBatcher struct {
	usecase *usecase.Usecase
}

func (t *tagBatcher) BatchGetTagsByRecordIDs(ctx context.Context, recordIDs []domain.RecordID) []*dataloader.Result[[]*domain.Tag] {
	results := make([]*dataloader.Result[[]*domain.Tag], len(recordIDs))

	indexs := make(map[domain.RecordID]int, len(recordIDs))
	for i, ID := range recordIDs {
		indexs[ID] = i
	}

	userID, err := ctxdef.UserID(ctx)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[[]*domain.Tag]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	tagWithRecordIDs, err := t.usecase.GetTagsWithRecordIDByRecordIDs(ctx, userID, recordIDs)
	if err != nil {
		for i := range results {
			results[i] = &dataloader.Result[[]*domain.Tag]{Error: xerrors.Errorf(": %w", err)}
		}
		return results
	}

	for _, recordID := range recordIDs {
		results[indexs[recordID]] = &dataloader.Result[[]*domain.Tag]{
			Data:  make([]*domain.Tag, 0),
			Error: nil,
		}
	}

	for _, tagWithRecordID := range tagWithRecordIDs {
		results[indexs[tagWithRecordID.RecordID]].Data = append(results[indexs[tagWithRecordID.RecordID]].Data, &tagWithRecordID.Tag)
	}

	return results
}
