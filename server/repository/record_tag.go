package repository

import (
	"context"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const recordTagTableName = "record_tag"

type RecordTagRepository struct {
	sess *dbr.Session
}

func NewRecordTagRepository(sess *dbr.Session) *RecordTagRepository {
	return &RecordTagRepository{
		sess: sess,
	}
}

func (r *RecordTagRepository) Insert(ctx context.Context, recordID domain.RecordID, tagID domain.TagID) error {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(recordTagTableName).
		Columns("record_id", "tag_id").
		Values(recordID, tagID).
		Exec()
	if err != nil {
		return xerrors.Errorf("failed to insert record_tag: %w", err)
	}

	return nil
}

func (r *RecordTagRepository) DeleteByRecordID(ctx context.Context, recordID domain.RecordID) error {
	runner := getRunner(ctx, r.sess)
	_, err := runner.DeleteFrom(recordTagTableName).
		Where("record_id = ?", recordID).
		Exec()
	if err != nil {
		return xerrors.Errorf("failed to delete record_tag by record_id: %w", err)
	}

	return nil
}
