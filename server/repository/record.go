package repository

import (
	"context"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const recordTableName = "record"

type RecordRepository struct {
	sess *dbr.Session
}

func NewRecordRepository(sess *dbr.Session) *RecordRepository {
	return &RecordRepository{
		sess: sess,
	}
}

func (r *RecordRepository) Insert(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(recordTableName).Columns("id", "user_id", "record_type", "title", "description", "at").Record(record).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert record: %w", err)
	}

	return record, nil
}

func (r *RecordRepository) GetMultiByUserIDAndAssetID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, assetID *domain.AssetID) ([]*domain.Record, *domain.PageInfo, error) {
	runner := getRunner(ctx, r.sess)
	records := make([]*domain.Record, 0)

	stmt := runner.Select("rc.*").From(dbr.I(recordTableName).As("rc")).Where("rc.user_id = ?", userID)

	if assetID != nil {
		stmt.Join(dbr.I(assetChangeTableName).As("ac"), "ac.record_id = rc.id").
			Where("ac.asset_id = ?", *assetID)
	}

	stmt, err := paginate(pageParam, stmt)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to paginate: %w", err)
	}

	_, err = stmt.LoadContext(ctx, &records)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to load records: %w", err)
	}

	var startCursor *domain.PageCursor
	var endCursor *domain.PageCursor
	if len(records) > 0 {
		switch pageParam.SortKey {
		case domain.RecordSortKeyAt.String():
			startCursor = domain.NewPageCursor(string(records[0].ID), records[0].At.Format("2006-01-02 15:04:05"))
			endCursor = domain.NewPageCursor(string(records[len(records)-1].ID), records[len(records)-1].At.Format("2006-01-02 15:04:05"))
		default:
			panic("unsupported sort key for record")
		}
	}

	hasNextPage, hasPreviousPage := hasPage(pageParam, len(records))

	pageInfo := &domain.PageInfo{
		StartCursor:     startCursor,
		EndCursor:       endCursor,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}

	return records, pageInfo, nil
}
