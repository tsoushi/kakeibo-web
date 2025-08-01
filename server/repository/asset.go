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
	_, err := runner.InsertInto(assettableName).Columns("id", "user_id", "name", "category_id").Record(asset).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert asset: %w", err)
	}

	return asset, nil
}

func (r *AssetRepository) Update(ctx context.Context, asset *domain.Asset) (*domain.Asset, error) {
	runner := getRunner(ctx, r.sess)
	result, err := runner.Update(assettableName).
		Set("name", asset.Name).
		Set("category_id", asset.CategoryID).
		Where("id = ? AND user_id = ?", asset.ID, asset.UserID).
		Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to update asset: %w", err)
	}
	resultCount, err := result.RowsAffected()
	if err != nil {
		return nil, xerrors.Errorf("failed to get affected rows: %w", err)
	}
	if resultCount == 0 {
		return nil, domain.ErrEntityNotFound
	}

	return asset, nil
}

func (r *AssetRepository) Delete(ctx context.Context, asset *domain.Asset) (domain.AssetID, error) {
	runner := getRunner(ctx, r.sess)
	result, err := runner.DeleteFrom(assettableName).
		Where("id = ? AND user_id = ?", asset.ID, asset.UserID).
		Exec()
	if err != nil {
		return "", xerrors.Errorf("failed to delete asset: %w", err)
	}
	resultCount, err := result.RowsAffected()
	if err != nil {
		return "", xerrors.Errorf("failed to get affected rows: %w", err)
	}
	if resultCount == 0 {
		return "", domain.ErrEntityNotFound
	}
	return asset.ID, nil
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

func (r *AssetRepository) GetMultiByUserIDAndCategoryID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, categoryID *domain.AssetCategoryID) ([]*domain.Asset, *domain.PageInfo, error) {
	runner := getRunner(ctx, r.sess)
	assets := make([]*domain.Asset, 0)

	stmt := runner.Select("*").From(assettableName).Where("user_id = ?", userID)

	if categoryID != nil {
		stmt = stmt.Where("category_id = ?", *categoryID)
	}

	stmt, err := paginate(pageParam, stmt)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to paginate: %w", err)
	}

	_, err = stmt.LoadContext(ctx, &assets)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to get assets by userID: %w", err)
	}

	var startCursor *domain.PageCursor
	var endCursor *domain.PageCursor
	if len(assets) > 0 {
		switch pageParam.SortKey {
		case domain.AssetSortKeyName.String():
			startCursor = domain.NewPageCursor(string(assets[0].ID), assets[0].Name)
			endCursor = domain.NewPageCursor(string(assets[len(assets)-1].ID), assets[len(assets)-1].Name)
		case domain.AssetCategorySortKeyCreatedAt.String():
			startCursor = domain.NewPageCursor(string(assets[0].ID), assets[0].CreatedAt.Format("2006-01-02 15:04:05"))
			endCursor = domain.NewPageCursor(string(assets[len(assets)-1].ID), assets[len(assets)-1].CreatedAt.Format("2006-01-02 15:04:05"))
		default:
			panic("unsupported sort key for asset repository")
		}
	}

	hasNextPage, hasPreviousPage := hasPage(pageParam, len(assets))

	pageInfo := &domain.PageInfo{
		StartCursor:     startCursor,
		EndCursor:       endCursor,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}

	return assets, pageInfo, nil
}

func (r *AssetRepository) GetMultiByUserIDAndIDs(ctx context.Context, userID domain.UserID, assetIDs []domain.AssetID) ([]*domain.Asset, error) {
	runner := getRunner(ctx, r.sess)
	assets := make([]*domain.Asset, 0)

	if len(assetIDs) == 0 {
		return assets, nil
	}

	_, err := runner.Select("*").
		From(assettableName).
		Where("user_id = ? AND id IN ?", userID, assetIDs).
		LoadContext(ctx, &assets)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to get assets by userID and IDs: %w", err)
	}

	return assets, nil
}

func (r *AssetRepository) GetMultiByCategoryIDs(ctx context.Context, userID domain.UserID, categoryIDs []domain.AssetCategoryID) ([]*domain.Asset, error) {
	runner := getRunner(ctx, r.sess)
	assets := make([]*domain.Asset, 0)

	if len(categoryIDs) == 0 {
		return assets, nil
	}

	_, err := runner.Select("*").
		From(assettableName).
		Where("user_id = ? AND category_id IN ?", userID, categoryIDs).
		LoadContext(ctx, &assets)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to get assets by userID and category IDs: %w", err)
	}

	return assets, nil
}
