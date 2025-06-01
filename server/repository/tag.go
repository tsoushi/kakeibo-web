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

func (r *TagRepository) GetMultiByUserID(ctx context.Context, pageParam *domain.PageParam, userID domain.UserID) ([]*domain.Tag, *domain.PageInfo, error) {
	runner := getRunner(ctx, r.sess)
	stmt := runner.Select("*").From(tagtableName).Where(dbr.Eq("user_id", userID))

	stmt, err := paginate(pageParam, stmt)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to paginate tags: %w", err)
	}

	tags := make([]*domain.Tag, 0)
	_, err = stmt.LoadContext(ctx, &tags)
	if err != nil {
		return nil, nil, xerrors.Errorf("failed to load tags: %w", err)
	}

	var startCursor *domain.PageCursor
	var endCursor *domain.PageCursor
	if len(tags) > 0 {
		if pageParam.SortKey == domain.TagSortKeyName.String() {
			startCursor = domain.NewPageCursor(string(tags[0].ID), tags[0].Name)
			endCursor = domain.NewPageCursor(string(tags[len(tags)-1].ID), tags[len(tags)-1].Name)
		} else {
			panic("unsupported sort key for asset repository")
		}
	}

	hasNextPage, hasPreviousPage := hasPage(pageParam, len(tags))

	pageInfo := &domain.PageInfo{
		StartCursor:     startCursor,
		EndCursor:       endCursor,
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
	}

	return tags, pageInfo, nil
}
