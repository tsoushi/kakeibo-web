package repository

import (
	"context"
	"errors"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const assettableName = "asset"

type AssetRepository struct {
	sess *dbr.Session
}

func NewAssetRepository(sess *dbr.Session) *AssetRepository {
	return &AssetRepository{
		sess: sess,
	}
}

func (r *AssetRepository) Insert(ctx context.Context, asset *domain.Asset) (*domain.Asset, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(assettableName).Columns("id", "user_id", "name").Record(asset).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert asset: %w", err)
	}

	return asset, nil
}

func (r *AssetRepository) List(ctx context.Context, userID domain.UserID) ([]*domain.Asset, error) {
	runner := getRunner(ctx, r.sess)
	assets := make([]*domain.Asset, 0)
	_, err := runner.Select("*").From(assettableName).Where("user_id = ?", userID).LoadContext(ctx, &assets)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to list assets by userID: %w", err)
	}

	return assets, nil
}
