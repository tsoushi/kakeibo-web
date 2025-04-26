package domain

import (
	"golang.org/x/xerrors"
)

var (
	ErrEntityNotFound = xerrors.New("entity not found")
)
