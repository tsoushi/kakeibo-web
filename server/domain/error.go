package domain

import (
	"golang.org/x/xerrors"
)

var (
	ErrEntityNotFound      = xerrors.New("entity not found")
	ErrInvalidPageParam    = xerrors.New("page param invalid")
	ErrInvalidRecordAmount = xerrors.New("invalid record amount")
)
