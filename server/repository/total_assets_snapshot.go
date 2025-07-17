package repository

import (
	"context"
	"errors"
	"kakeibo-web-server/domain"
	"time"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const totalAssetsSnapshotTableName = "total_assets_snapshot"

type TotalAssetsSnapshotRepository struct {
	sess *dbr.Session
}

func NewTotalAssetsSnapshotRepository(sess *dbr.Session) *TotalAssetsSnapshotRepository {
	return &TotalAssetsSnapshotRepository{
		sess: sess,
	}
}

func (r *TotalAssetsSnapshotRepository) Insert(ctx context.Context, snapshot *domain.TotalAssetsSnapshot) (*domain.TotalAssetsSnapshot, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(totalAssetsSnapshotTableName).
		Columns("id", "user_id", "asset_id", "amount", "at", "is_valid").
		Record(snapshot).
		Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert total assets snapshot: %w", err)
	}

	return snapshot, nil
}

func (r *TotalAssetsSnapshotRepository) OptionalGetValidLatestByUserIDAndAssetIDAndBefore(ctx context.Context, userID domain.UserID, assetID *domain.AssetID, before time.Time) (*domain.TotalAssetsSnapshot, error) {
	runner := getRunner(ctx, r.sess)
	snapshot := &domain.TotalAssetsSnapshot{}

	stmt := runner.Select("*").
		From(totalAssetsSnapshotTableName).
		Where("user_id = ?", userID).
		Where("is_valid = ?", true).
		Where("at < ?", before)

	if assetID != nil {
		stmt = stmt.Where("asset_id = ?", *assetID)
	}

	stmt = stmt.OrderDir("at", false).Limit(1)

	err := stmt.LoadOneContext(ctx, snapshot)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, nil
		}
		return nil, xerrors.Errorf("failed to get valid recent total assets snapshot: %w", err)
	}

	return snapshot, nil
}
