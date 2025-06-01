package repository

import (
	"context"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const tagtableName = "tag"

type TagRepository struct {
	sess *dbr.Session
}

func NewTagRepository(sess *dbr.Session) *TagRepository {
	return &TagRepository{
		sess: sess,
	}
}

func (r *TagRepository) Insert(cxt context.Context, tag *domain.Tag) (*domain.Tag, error) {
	runner := getRunner(cxt, r.sess)
	_, err := runner.InsertInto(tagtableName).
		Columns("id", "user_id", "name").
		Record(tag).
		Exec()

	if err != nil {
		return nil, xerrors.Errorf("failed to insert tag: %w", err)
	}

	return tag, nil
}
