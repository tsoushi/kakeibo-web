package usecase

import (
	"context"
	"kakeibo-web-server/domain"
	"time"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateIncomeRecord(ctx context.Context, userID domain.UserID, title string, description string, at time.Time, assetID domain.AssetID, amount int) (*domain.Record, *domain.AssetChange, error) {
	record, assetChange, err := domain.NewRecordIncomeWithAssetChange(userID, title, description, at, assetID, amount)
	if err != nil {
		return nil, nil, err
	}

	err = u.repo.RunInTx(ctx, func(ctx context.Context) error {
		_, err = u.repo.Record.Insert(ctx, record)
		if err != nil {
			return xerrors.Errorf("failed to insert record: %w", err)
		}

		_, err = u.repo.AssetChange.Insert(ctx, assetChange)
		if err != nil {
			return xerrors.Errorf("failed to insert asset change: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return record, assetChange, nil
}

func (u *Usecase) CreateExpenseRecord(ctx context.Context, userID domain.UserID, title string, description string, at time.Time, assetID domain.AssetID, amount int) (*domain.Record, *domain.AssetChange, error) {
	record, assetChange, err := domain.NewRecordExpenseWithAssetChange(userID, title, description, at, assetID, amount)
	if err != nil {
		return nil, nil, err
	}

	err = u.repo.RunInTx(ctx, func(ctx context.Context) error {
		_, err = u.repo.Record.Insert(ctx, record)
		if err != nil {
			return xerrors.Errorf("failed to insert record: %w", err)
		}

		_, err = u.repo.AssetChange.Insert(ctx, assetChange)
		if err != nil {
			return xerrors.Errorf("failed to insert asset change: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return record, assetChange, nil
}

func (u *Usecase) CreateTransferRecord(ctx context.Context, userID domain.UserID, title string, description string, at time.Time, fromAssetID domain.AssetID, toAssetID domain.AssetID, amount int) (*domain.Record, *domain.AssetChange, *domain.AssetChange, error) {
	record, fromAssetChange, toAssetChange, err := domain.NewRecordTransferWithAssetChanges(userID, title, description, at, fromAssetID, toAssetID, amount)
	if err != nil {
		return nil, nil, nil, xerrors.Errorf("failed to create transfer record: %w", err)
	}

	err = u.repo.RunInTx(ctx, func(ctx context.Context) error {
		_, err = u.repo.Record.Insert(ctx, record)
		if err != nil {
			return xerrors.Errorf("failed to insert record: %w", err)
		}

		_, err = u.repo.AssetChange.Insert(ctx, fromAssetChange)
		if err != nil {
			return xerrors.Errorf("failed to insert from asset change: %w", err)
		}
		_, err = u.repo.AssetChange.Insert(ctx, toAssetChange)
		if err != nil {
			return xerrors.Errorf("failed to insert to asset change: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}

	return record, fromAssetChange, toAssetChange, nil
}

func (u *Usecase) GetRecordsByUserIDAndAssetID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, assetID *domain.AssetID) (domain.Records, *domain.PageInfo, error) {
	records, pageInfo, err := u.repo.Record.GetMultiByUserIDAndAssetID(ctx, pageParam, userID, assetID)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to get records: %w", err)
	}

	return records, pageInfo, nil
}

func (u *Usecase) CulcTotalAssetAmountAndCreateSnapshot(ctx context.Context, userID domain.UserID, assetID *domain.AssetID, beforeRecord domain.Record) (int, error) {
	totalAssetsSnapshot, err := u.repo.TotalAssetsSnapshot.OptionalGetValidLatestByUserIDAndAssetIDAndBefore(ctx, userID, assetID, beforeRecord.At)
	if err != nil {
		return 0, xerrors.Errorf("failed to get total assets snapshot: %w", err)
	}

	totalAssetsAmount := 0
	totalAssetsAt := time.Time{}
	if totalAssetsSnapshot != nil {
		totalAssetsAmount = totalAssetsSnapshot.Amount
		totalAssetsAt = totalAssetsSnapshot.At
	}

	assetChangeWithAts, err := u.repo.AssetChange.GetMultiWithAtByAssetIDAndAfterAndBefore(ctx, userID, assetID, totalAssetsAt, beforeRecord.At, beforeRecord.ID)
	if err != nil {
		return 0, xerrors.Errorf("failed to get asset changes: %w", err)
	}

	if len(assetChangeWithAts) > domain.TotalAssetsSnapshotDefaultSpan {
		newTotalAssetsSnapshots := assetChangeWithAts.CreateSnapshots(assetID, totalAssetsAmount, domain.TotalAssetsSnapshotDefaultSpan)

		for _, snapshot := range newTotalAssetsSnapshots {
			_, err := u.repo.TotalAssetsSnapshot.Insert(ctx, snapshot)
			if err != nil {
				return 0, xerrors.Errorf("failed to insert total assets snapshot: %w", err)
			}
		}
	}

	totalAssetsAmount += assetChangeWithAts.TotalAmount()

	return totalAssetsAmount, nil
}
