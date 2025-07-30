package usecase

import (
	"context"
	"kakeibo-web-server/domain"
	"time"

	"golang.org/x/xerrors"
)

func (u *Usecase) CreateIncomeRecord(ctx context.Context, userID domain.UserID, title string, description string, at time.Time, assetID domain.AssetID, amount int, tagNames []string) (*domain.Record, *domain.AssetChange, error) {
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

		tags, err := u.GetOrCreateTagsByName(ctx, userID, tagNames)
		if err != nil {
			return xerrors.Errorf("failed to get or create tags: %w", err)
		}

		for _, tag := range tags {
			err = u.repo.RecordTag.Insert(ctx, record.ID, tag.ID)
			if err != nil {
				return xerrors.Errorf("failed to insert record tag: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return record, assetChange, nil
}

func (u *Usecase) CreateExpenseRecord(ctx context.Context, userID domain.UserID, title string, description string, at time.Time, assetID domain.AssetID, amount int, tagNames []string) (*domain.Record, *domain.AssetChange, error) {
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

		tags, err := u.GetOrCreateTagsByName(ctx, userID, tagNames)
		if err != nil {
			return xerrors.Errorf("failed to get or create tags: %w", err)
		}
		for _, tag := range tags {
			err = u.repo.RecordTag.Insert(ctx, record.ID, tag.ID)
			if err != nil {
				return xerrors.Errorf("failed to insert record tag: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, xerrors.Errorf(": %w", err)
	}

	return record, assetChange, nil
}

func (u *Usecase) CreateTransferRecord(ctx context.Context, userID domain.UserID, title string, description string, at time.Time, fromAssetID domain.AssetID, toAssetID domain.AssetID, amount int, tagNames []string) (*domain.Record, *domain.AssetChange, *domain.AssetChange, error) {
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

		tags, err := u.GetOrCreateTagsByName(ctx, userID, tagNames)
		if err != nil {
			return xerrors.Errorf("failed to get or create tags: %w", err)
		}
		for _, tag := range tags {
			err = u.repo.RecordTag.Insert(ctx, record.ID, tag.ID)
			if err != nil {
				return xerrors.Errorf("failed to insert record tag: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return nil, nil, nil, xerrors.Errorf(": %w", err)
	}

	return record, fromAssetChange, toAssetChange, nil
}

func (u *Usecase) UpdateIncomeRecord(ctx context.Context, userID domain.UserID, id domain.RecordID, title string, description string, at time.Time, assetID domain.AssetID, amount int, tagNames []string) (*domain.Record, error) {
	var record *domain.Record
	err := u.repo.RunInTx(ctx, func(ctx context.Context) error {
		getRecord, err := u.repo.Record.GetByID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to get record by ID: %w", err)
		}

		record = getRecord

		record.Title = title
		record.Description = description
		record.At = at

		_, err = u.repo.Record.Update(ctx, record)
		if err != nil {
			return xerrors.Errorf("failed to update record: %w", err)
		}

		assetChanges, err := u.repo.AssetChange.GetMultiByRecordID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to get asset changes by record ID: %w", err)
		}
		if len(assetChanges) == 0 {
			return xerrors.Errorf("asset change not found for record ID: %w", domain.ErrEntityNotFound)
		} else if len(assetChanges) > 1 {
			return xerrors.Errorf("multiple asset changes found for record ID: %w", domain.ErrEntityNotFound)
		}

		incomeAssetChange := assetChanges.Income()

		incomeAssetChange.AssetID = assetID
		incomeAssetChange.Amount = amount
		_, err = u.repo.AssetChange.Update(ctx, incomeAssetChange)
		if err != nil {
			return xerrors.Errorf("failed to update asset change: %w", err)
		}

		tags, err := u.GetOrCreateTagsByName(ctx, userID, tagNames)
		if err != nil {
			return xerrors.Errorf("failed to get or create tags: %w", err)
		}
		err = u.repo.RecordTag.DeleteByRecordID(ctx, record.ID)
		if err != nil {
			return xerrors.Errorf("failed to delete record tags: %w", err)
		}
		for _, tag := range tags {
			err = u.repo.RecordTag.Insert(ctx, record.ID, tag.ID)
			if err != nil {
				return xerrors.Errorf("failed to insert record tag: %w", err)
			}
		}

		u.repo.TotalAssetsSnapshot.InvalidateByUserIDAndSince(ctx, userID, record.At)

		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return record, nil
}

func (u *Usecase) UpdateExpenseRecord(ctx context.Context, userID domain.UserID, id domain.RecordID, title string, description string, at time.Time, assetID domain.AssetID, amount int, tagNames []string) (*domain.Record, error) {
	var record *domain.Record
	err := u.repo.RunInTx(ctx, func(ctx context.Context) error {
		getRecord, err := u.repo.Record.GetByID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to get record by ID: %w", err)
		}

		record = getRecord

		record.Title = title
		record.Description = description
		record.At = at

		_, err = u.repo.Record.Update(ctx, record)
		if err != nil {
			return xerrors.Errorf("failed to update record: %w", err)
		}

		assetChanges, err := u.repo.AssetChange.GetMultiByRecordID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to get asset changes by record ID: %w", err)
		}
		if len(assetChanges) == 0 {
			return xerrors.Errorf("asset change not found for record ID: %w", domain.ErrEntityNotFound)
		} else if len(assetChanges) > 1 {
			return xerrors.Errorf("multiple asset changes found for record ID: %w", domain.ErrEntityNotFound)
		}

		expenseAssetChange := assetChanges.Expense()
		expenseAssetChange.AssetID = assetID
		expenseAssetChange.Amount = -amount
		_, err = u.repo.AssetChange.Update(ctx, expenseAssetChange)
		if err != nil {
			return xerrors.Errorf("failed to update asset change: %w", err)
		}

		tags, err := u.GetOrCreateTagsByName(ctx, userID, tagNames)
		if err != nil {
			return xerrors.Errorf("failed to get or create tags: %w", err)
		}
		err = u.repo.RecordTag.DeleteByRecordID(ctx, record.ID)
		if err != nil {
			return xerrors.Errorf("failed to delete record tags: %w", err)
		}
		for _, tag := range tags {
			err = u.repo.RecordTag.Insert(ctx, record.ID, tag.ID)
			if err != nil {
				return xerrors.Errorf("failed to insert record tag: %w", err)
			}
		}

		u.repo.TotalAssetsSnapshot.InvalidateByUserIDAndSince(ctx, userID, record.At)

		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return record, nil
}

func (u *Usecase) UpdateTransferRecord(ctx context.Context, userID domain.UserID, id domain.RecordID, title string, description string, at time.Time, fromAssetID domain.AssetID, toAssetID domain.AssetID, amount int, tagNames []string) (*domain.Record, error) {
	var record *domain.Record
	err := u.repo.RunInTx(ctx, func(ctx context.Context) error {
		getRecord, err := u.repo.Record.GetByID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to get record by ID: %w", err)
		}

		record = getRecord

		record.Title = title
		record.Description = description
		record.At = at

		_, err = u.repo.Record.Update(ctx, record)
		if err != nil {
			return xerrors.Errorf("failed to update record: %w", err)
		}

		assetChanges, err := u.repo.AssetChange.GetMultiByRecordID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to get asset changes by record ID: %w", err)
		}
		if len(assetChanges) != 2 {
			return xerrors.Errorf("expected 2 asset changes for transfer record, found %d: %w", len(assetChanges), domain.ErrEntityNotFound)
		}
		fromAssetChange := assetChanges.Expense()
		toAssetChange := assetChanges.Income()
		fromAssetChange.AssetID = fromAssetID
		fromAssetChange.Amount = -amount
		toAssetChange.AssetID = toAssetID
		toAssetChange.Amount = amount
		_, err = u.repo.AssetChange.Update(ctx, fromAssetChange)
		if err != nil {
			return xerrors.Errorf("failed to update from asset change: %w", err)
		}
		_, err = u.repo.AssetChange.Update(ctx, toAssetChange)
		if err != nil {
			return xerrors.Errorf("failed to update to asset change: %w", err)
		}

		tags, err := u.GetOrCreateTagsByName(ctx, userID, tagNames)
		if err != nil {
			return xerrors.Errorf("failed to get or create tags: %w", err)
		}
		err = u.repo.RecordTag.DeleteByRecordID(ctx, record.ID)
		if err != nil {
			return xerrors.Errorf("failed to delete record tags: %w", err)
		}
		for _, tag := range tags {
			err = u.repo.RecordTag.Insert(ctx, record.ID, tag.ID)
			if err != nil {
				return xerrors.Errorf("failed to insert record tag: %w", err)
			}
		}

		u.repo.TotalAssetsSnapshot.InvalidateByUserIDAndSince(ctx, userID, record.At)

		return nil
	})
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return record, nil
}

func (u *Usecase) DeleteRecord(ctx context.Context, userID domain.UserID, id domain.RecordID) (domain.RecordID, error) {
	err := u.repo.RunInTx(ctx, func(ctx context.Context) error {
		record, err := u.repo.Record.GetByID(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("record not found: %w", err)
		}

		_, err = u.repo.Record.Delete(ctx, userID, id)
		if err != nil {
			return xerrors.Errorf("failed to delete record: %w", err)
		}

		u.repo.TotalAssetsSnapshot.InvalidateByUserIDAndSince(ctx, userID, record.At)

		return nil
	})
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}

	return id, nil
}

func (u *Usecase) GetRecordByID(ctx context.Context, userID domain.UserID, id domain.RecordID) (*domain.Record, error) {
	record, err := u.repo.Record.GetByID(ctx, userID, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to get record by ID: %w", err)
	}

	return record, nil
}

func (u *Usecase) GetRecordsByUserIDAndAssetID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, assetID *domain.AssetID) (domain.Records, *domain.PageInfo, error) {
	records, pageInfo, err := u.repo.Record.GetMultiByUserIDAndAssetID(ctx, pageParam, userID, assetID)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to get records: %w", err)
	}

	return records, pageInfo, nil
}

func (u *Usecase) CulcTotalAssetAmountAndCreateSnapshot(ctx context.Context, userID domain.UserID, assetID *domain.AssetID, before time.Time, recordID domain.RecordID) (int, error) {
	totalAssetsSnapshot, err := u.repo.TotalAssetsSnapshot.OptionalGetValidLatestByUserIDAndAssetIDAndBefore(ctx, userID, assetID, before)
	if err != nil {
		return 0, xerrors.Errorf("failed to get total assets snapshot: %w", err)
	}

	totalAssetsAmount := 0
	totalAssetsAt := time.Time{}
	if totalAssetsSnapshot != nil {
		totalAssetsAmount = totalAssetsSnapshot.Amount
		totalAssetsAt = totalAssetsSnapshot.At
	}

	assetChangeWithAts, err := u.repo.AssetChange.GetMultiWithAtByAssetIDAndAfterAndBefore(ctx, userID, assetID, totalAssetsAt, before, recordID)
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

func (u *Usecase) GetRecordsPerMonth(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, year int, month int, tagNames []string, assetIDs []domain.AssetID, recordTypes []domain.RecordType) (domain.Records, *domain.PageInfo, error) {
	records, pageInfo, err := u.repo.Record.GetMultiByUserIDAndYearAndMonth(ctx, pageParam, userID, year, month, tagNames, assetIDs, recordTypes)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to get records per month: %w", err)
	}

	return records, pageInfo, nil
}
