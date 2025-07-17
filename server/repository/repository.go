package repository

import (
	"context"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
)

type TxKey struct{}

type Repository struct {
	sess                *dbr.Session
	User                *UserRepository
	Asset               *AssetRepository
	AssetCategory       *AssetCategoryRepository
	Tag                 *TagRepository
	Record              *RecordRepository
	AssetChange         *AssetChangeRepository
	TotalAssetsSnapshot *TotalAssetsSnapshotRepository
}

func NewRepository(sess *dbr.Session) *Repository {
	return &Repository{
		sess:                sess,
		User:                NewUserRepository(sess),
		Asset:               NewAssetRepository(sess),
		AssetCategory:       NewAssetCategoryRepository(sess),
		Tag:                 NewTagRepository(sess),
		Record:              NewRecordRepository(sess),
		AssetChange:         NewAssetChangeRepository(sess),
		TotalAssetsSnapshot: NewTotalAssetsSnapshotRepository(sess),
	}
}

func (r *Repository) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := r.sess.Begin()
	if err != nil {
		return err
	}
	defer tx.RollbackUnlessCommitted()

	ctx = context.WithValue(ctx, TxKey{}, tx)

	err = fn(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getRunner(ctx context.Context, sess *dbr.Session) dbr.SessionRunner {
	value := ctx.Value(TxKey{})
	if value == nil {
		// RunInTxの中でない場合は、sessを返す
		return sess
	}

	tx, ok := value.(*dbr.Tx)
	if !ok {
		panic("context value is not a transaction")
	}

	return tx
}

func paginate(pageParam *domain.PageParam, stmt *dbr.SelectStmt) (*dbr.SelectStmt, error) {
	if pageParam == nil {
		return stmt, nil
	}

	if pageParam.First != nil && pageParam.Last != nil {
		return nil, domain.ErrInvalidPageParam
	}
	if pageParam.First != nil {
		stmt.Limit(uint64(*pageParam.First))
		if pageParam.After != nil {
			stmt.Where(
				dbr.Or(
					dbr.Gt(pageParam.SortKey, pageParam.After.Value),
					dbr.And(
						dbr.Eq(pageParam.SortKey, pageParam.After.Value),
						dbr.Gt("id", pageParam.After.ID),
					),
				),
			)
		}
	}

	if pageParam.Last != nil {
		stmt.Limit(uint64(*pageParam.Last))
		if pageParam.Before != nil {
			stmt.Where(
				dbr.Or(
					dbr.Lt(pageParam.SortKey, pageParam.Before.Value),
					dbr.And(
						dbr.Eq(pageParam.SortKey, pageParam.Before.Value),
						dbr.Lt("id", pageParam.Before.ID),
					),
				),
			)
		}
	}

	if pageParam.IsForward() {
		stmt = stmt.OrderAsc(pageParam.SortKey)
		stmt = stmt.OrderAsc("id")
	} else {
		stmt = stmt.OrderDesc(pageParam.SortKey)
		stmt = stmt.OrderDesc("id")
	}

	return stmt, nil
}

func hasPage(pageParam *domain.PageParam, resultCount int) (hasNextPage, hasPreviousPage bool) {
	hasNextPage = false
	hasPreviousPage = false

	if pageParam.After != nil {
		hasPreviousPage = true
	} else if pageParam.Before != nil {
		hasNextPage = true
	}

	if pageParam.First != nil && resultCount == *pageParam.First {
		hasNextPage = true
	} else if pageParam.Last != nil && resultCount == *pageParam.Last {
		hasPreviousPage = true
	}

	return hasNextPage, hasPreviousPage
}
