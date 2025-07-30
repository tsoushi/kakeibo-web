package repository

import (
	"context"
	"errors"
	"kakeibo-web-server/domain"
	"time"

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

func (r *RecordRepository) GetByID(ctx context.Context, userID domain.UserID, id domain.RecordID) (*domain.Record, error) {
	runner := getRunner(ctx, r.sess)
	record := &domain.Record{}
	err := runner.Select("*").From(recordTableName).Where("id = ? AND user_id = ?", id, userID).LoadOne(record)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to get record by ID: %w", err)
	}
	return record, nil
}

func (r *RecordRepository) Update(ctx context.Context, record *domain.Record) (*domain.Record, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.Update(recordTableName).
		Set("title", record.Title).
		Set("description", record.Description).
		Set("at", record.At).
		Where("id = ? AND user_id = ?", record.ID, record.UserID).
		Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to update record: %w", err)
	}

	return record, nil
}

func (r *RecordRepository) Delete(ctx context.Context, userID domain.UserID, id domain.RecordID) (*domain.Record, error) {
	runner := getRunner(ctx, r.sess)
	record, err := r.GetByID(ctx, userID, id)
	if err != nil {
		return nil, xerrors.Errorf("failed to get record for deletion: %w", err)
	}
	_, err = runner.DeleteFrom(recordTableName).Where("id = ? AND user_id = ?", id, userID).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to delete record: %w", err)
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

func (r *RecordRepository) GetMultiByUserIDAndYearAndMonth(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID, year int, month int, tagNames []string, assetIDs []domain.AssetID, recordTypes []domain.RecordType) ([]*domain.Record, *domain.PageInfo, error) {
	runner := getRunner(ctx, r.sess)
	records := make([]*domain.Record, 0)

	stmt := runner.Select("rc.*").From(dbr.I(recordTableName).As("rc")).Where("rc.user_id = ?", userID)

	if len(tagNames) > 0 {
		stmt.Join(dbr.I("record_tag").As("rt"), "rt.record_id = rc.id").
			Join(dbr.I("tag").As("t"), "t.id = rt.tag_id").
			Where("t.name IN ?", tagNames)
	}

	if len(assetIDs) > 0 {
		stmt.Join(dbr.I(assetChangeTableName).As("ac"), "ac.record_id = rc.id").
			Where("ac.asset_id IN ?", assetIDs)
	}

	if len(recordTypes) > 0 {
		stmt.Where("rc.record_type IN ?", recordTypes)
	}

	stmt.Where(dbr.And(
		dbr.Gte("rc.at", time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)),
		dbr.Lt("rc.at", time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, time.Local)),
	))

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
