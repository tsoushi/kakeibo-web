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

func (r *TagRepository) Update(ctx context.Context, tag *domain.Tag) (*domain.Tag, error) {
	runner := getRunner(ctx, r.sess)
	result, err := runner.Update(tagtableName).
		Set("name", tag.Name).
		Where("id = ? AND user_id = ?", tag.ID, tag.UserID).
		Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to update tag: %w", err)
	}
	resultCount, err := result.RowsAffected()
	if err != nil {
		return nil, xerrors.Errorf("failed to get affected rows: %w", err)
	}
	if resultCount == 0 {
		return nil, domain.ErrEntityNotFound
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

func (r *TagRepository) Delete(ctx context.Context, userID domain.UserID, tagID domain.TagID) (domain.TagID, error) {
	runner := getRunner(ctx, r.sess)

	result, err := runner.DeleteFrom(tagtableName).
		Where(dbr.Eq("id", tagID)).
		Where(dbr.Eq("user_id", userID)).
		Exec()

	if err != nil {
		return "", xerrors.Errorf("failed to delete tag: %w", err)
	}

	resultCount, err := result.RowsAffected()
	if err != nil {
		return "", xerrors.Errorf("failed to get affected rows: %w", err)
	}
	if resultCount == 0 {
		return "", domain.ErrEntityNotFound
	}

	return tagID, nil
}

func (r *TagRepository) GetMultiByNames(ctx context.Context, userID domain.UserID, names []string) (domain.Tags, error) {
	if len(names) == 0 {
		return nil, nil
	}

	runner := getRunner(ctx, r.sess)
	tags := make([]*domain.Tag, 0, len(names))

	_, err := runner.Select("*").From(tagtableName).
		Where(dbr.Eq("user_id", userID)).
		Where("name IN ?", names).
		LoadContext(ctx, &tags)

	if err != nil {
		return nil, xerrors.Errorf("failed to load tags by names: %w", err)
	}

	return tags, nil
}

func (r *TagRepository) GetMultiWithRecordIDByRecordIDs(ctx context.Context, userID domain.UserID, recordIDs []domain.RecordID) ([]*domain.TagWithRecordID, error) {
	if len(recordIDs) == 0 {
		return nil, nil
	}

	runner := getRunner(ctx, r.sess)
	tagWithRecordIDs := make([]*domain.TagWithRecordID, 0)

	_, err := runner.Select("rt.record_id, t.*").From(dbr.I(recordTagTableName).As("rt")).
		Join(dbr.I(tagtableName).As("t"), "t.id = rt.tag_id").
		Where(dbr.Eq("t.user_id", userID)).
		Where("rt.record_id IN ?", recordIDs).
		LoadContext(ctx, &tagWithRecordIDs)

	if err != nil {
		return nil, xerrors.Errorf("failed to load tags by record IDs: %w", err)
	}

	return tagWithRecordIDs, nil
}
