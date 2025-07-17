package usecase

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

func (u *Usecase) GetAssetChangesByRecordIDs(ctx context.Context, userID domain.UserID, recordIDs []domain.RecordID) ([]*domain.AssetChange, error) {
	assetChanges, err := u.repo.AssetChange.GetMultiByRecordIDs(ctx, userID, recordIDs)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return assetChanges, nil
}
