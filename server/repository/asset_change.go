package repository

import (
	"context"
	"kakeibo-web-server/domain"
	"time"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const assetChangeTableName = "asset_change"

type AssetChangeRepository struct {
	sess *dbr.Session
}

func NewAssetChangeRepository(sess *dbr.Session) *AssetChangeRepository {
	return &AssetChangeRepository{
		sess: sess,
	}
}

func (r *AssetChangeRepository) Insert(ctx context.Context, change *domain.AssetChange) (*domain.AssetChange, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(assetChangeTableName).
		Columns("id", "user_id", "record_id", "asset_id", "amount").
		Record(change).
		Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert asset change: %w", err)
	}

	return change, nil
}

func (r *AssetChangeRepository) GetMultiByRecordIDs(ctx context.Context, userID domain.UserID, recordIDs []domain.RecordID) ([]*domain.AssetChange, error) {
	runner := getRunner(ctx, r.sess)
	changes := make([]*domain.AssetChange, 0)

	if len(recordIDs) == 0 {
		return changes, nil
	}

	_, err := runner.Select("*").
		From(assetChangeTableName).
		Where("user_id = ? AND record_id IN ?", userID, recordIDs).
		LoadContext(ctx, &changes)
	if err != nil {
		return nil, xerrors.Errorf("failed to get asset changes by record IDs: %w", err)
	}

	return changes, nil
}

func (r *AssetChangeRepository) GetMultiWithAtByAssetIDAndAfterAndBefore(ctx context.Context, userID domain.UserID, assetID *domain.AssetID, after, before time.Time, beforeRecordID domain.RecordID) (domain.AssetChangeWithAts, error) {
	runner := getRunner(ctx, r.sess)
	changes := make([]*domain.AssetChangeWithAt, 0)

	stmt := runner.Select("ac.*, rc.at").
		From(dbr.I(assetChangeTableName).As("ac")).
		Where("ac.user_id = ?", userID).
		LeftJoin(dbr.I("record").As("rc"), "rc.id = ac.record_id").
		Where(dbr.Or(
			dbr.And(dbr.Gt("rc.at", after), dbr.Lt("rc.at", before)),
			dbr.And(dbr.Eq("rc.at", before), dbr.Lt("ac.record_id", beforeRecordID)),
		)).
		OrderAsc("rc.at").
		OrderAsc("ac.id")

	if assetID != nil {
		stmt = stmt.Where("ac.asset_id = ?", *assetID)
	}

	_, err := stmt.LoadContext(ctx, &changes)
	if err != nil {
		return nil, xerrors.Errorf("failed to get asset changes by asset ID and time range: %w", err)
	}

	return changes, nil
}
