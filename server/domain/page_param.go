package domain

import (
	"golang.org/x/xerrors"
)

type PageParam struct {
	First   *int
	After   *PageCursor
	Last    *int
	Before  *PageCursor
	SortKey string
}

// NewPageParam でフィールドの整合性チェックを行うため、New後にフィールドを直接変更しないこと（panicになる）
func NewPageParam(first *int, after *PageCursor, last *int, before *PageCursor, sortKey string) (*PageParam, error) {
	if first != nil && last != nil {
		return nil, xerrors.Errorf("first and last cannot be set at the same time: %w", ErrInvalidPageParam)
	}
	if first == nil && last == nil {
		return nil, xerrors.Errorf("either first or last must be set: %w", ErrInvalidPageParam)
	}
	if first != nil && before != nil {
		return nil, xerrors.Errorf("first and before cannot be set at the same time: %w", ErrInvalidPageParam)
	}
	if last != nil && after != nil {
		return nil, xerrors.Errorf("last and after cannot be set at the same time: %w", ErrInvalidPageParam)
	}

	return &PageParam{
		First:   first,
		After:   after,
		Last:    last,
		Before:  before,
		SortKey: sortKey,
	}, nil
}

func (p PageParam) IsReverse() bool {
	if p.First != nil && p.Last == nil {
		return false
	} else if p.First == nil && p.Last != nil {
		return true
	}

	panic("invalid PageParam")
}
