package repository

import (
	"context"
	"errors"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

type AssetCategoryRepository struct {
	sess *dbr.Session
}

func NewAssetCategoryRepository(sess *dbr.Session) *AssetCategoryRepository {
	return &AssetCategoryRepository{
		sess: sess,
	}
}

func (r *AssetCategoryRepository) Insert(ctx context.Context, category *domain.AssetCategory) (*domain.AssetCategory, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto("asset_category").Columns("id", "user_id", "name").Record(category).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert asset category: %w", err)
	}

	return category, nil
}

func (r *AssetCategoryRepository) Update(ctx context.Context, category *domain.AssetCategory) (*domain.AssetCategory, error) {
	runner := getRunner(ctx, r.sess)
	result, err := runner.Update("asset_category").
		Set("name", category.Name).
		Where("id = ? AND user_id = ?", category.ID, category.UserID).
		Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to update asset category: %w", err)
	}

	resultCount, err := result.RowsAffected()
	if err != nil {
		return nil, xerrors.Errorf("failed to get affected rows: %w", err)
	}
	if resultCount == 0 {
		return nil, domain.ErrEntityNotFound
	}

	return category, nil
}

func (r *AssetCategoryRepository) GetByID(ctx context.Context, userID domain.UserID, id domain.AssetCategoryID) (*domain.AssetCategory, error) {
	runner := getRunner(ctx, r.sess)
	category := &domain.AssetCategory{}

	_, err := runner.Select("*").From("asset_category").Where("id = ? AND user_id = ?", id, userID).LoadContext(ctx, category)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to get asset category by ID: %w", err)
	}

	return category, nil
}

func (r *AssetCategoryRepository) GetMultiByUserID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID) ([]*domain.AssetCategory, *domain.PageInfo, error) {
	runner := getRunner(ctx, r.sess)
	categories := make([]*domain.AssetCategory, 0)

	stmt := runner.Select("*").From("asset_category").Where("user_id = ?", userID)

	stmt, err := paginate(pageParam, stmt)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to paginate: %w", err)
	}

	_, err = stmt.LoadContext(ctx, &categories)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, nil, domain.ErrEntityNotFound
		}
		return nil, nil, xerrors.Errorf("failed to get asset categories by userID: %w", err)
	}

	var startCursor *domain.PageCursor
	var endCursor *domain.PageCursor
	if len(categories) > 0 {
		switch pageParam.SortKey {
		case domain.AssetCategorySortKeyName.String():
			startCursor = domain.NewPageCursor(string(categories[0].ID), categories[0].Name)
			endCursor = domain.NewPageCursor(string(categories[len(categories)-1].ID), categories[len(categories)-1].Name)
		case domain.AssetCategorySortKeyCreatedAt.String():
			startCursor = domain.NewPageCursor(string(categories[0].ID), categories[0].CreatedAt.Format("2006-01-02 15:04:05"))
			endCursor = domain.NewPageCursor(string(categories[len(categories)-1].ID), categories[len(categories)-1].CreatedAt.Format("2006-01-02 15:04:05"))
		default:
			panic("unknown sort key")
		}
	}

	hasNextPage, hasPreviousPage := hasPage(pageParam, len(categories))

	pageInfo := &domain.PageInfo{
		StartCursor:     startCursor,
		EndCursor:       endCursor,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}

	return categories, pageInfo, nil
}

func (r *AssetCategoryRepository) GetMultiByAssetCategoryIDs(ctx context.Context, userID domain.UserID, assetCategoryIDs []domain.AssetCategoryID) ([]*domain.AssetCategory, error) {
	categories := make([]*domain.AssetCategory, 0)

	stmt := r.sess.Select("*").From("asset_category").
		Where("user_id = ?", userID).
		Where("id IN ?", assetCategoryIDs)

	_, err := stmt.LoadContext(ctx, &categories)
	if err != nil {
		return nil, xerrors.Errorf("failed to get asset categories by ID: %w", err)
	}

	return categories, nil
}

func (r *AssetCategoryRepository) Delete(ctx context.Context, userID domain.UserID, assetCategoryID domain.AssetCategoryID) (domain.AssetCategoryID, error) {
	runner := getRunner(ctx, r.sess)

	result, err := runner.DeleteFrom("asset_category").
		Where("id = ? AND user_id = ?", assetCategoryID, userID).
		Exec()
	if err != nil {
		return "", xerrors.Errorf("failed to delete asset category: %w", err)
	}

	resultCount, err := result.RowsAffected()
	if err != nil {
		return "", xerrors.Errorf("failed to get affected rows: %w", err)
	}
	if resultCount == 0 {
		return "", domain.ErrEntityNotFound
	}

	return assetCategoryID, nil
}
