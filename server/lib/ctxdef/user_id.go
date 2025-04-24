package ctxdef

import (
	"context"
	"kakeibo-web-server/domain"

	"golang.org/x/xerrors"
)

type UserIDKey struct{}

func WithUserID(ctx context.Context, userID domain.UserID) context.Context {
	return context.WithValue(ctx, UserIDKey{}, userID)
}

func MustUserID(ctx context.Context) domain.UserID {
	value := ctx.Value(UserIDKey{})
	if value == nil {
		panic("userID not found in context")
	}

	userID := value.(domain.UserID)

	return userID
}

func UserID(ctx context.Context) (domain.UserID, error) {
	value := ctx.Value(UserIDKey{})
	if value == nil {
		return "", xerrors.New("userID not found in context")
	}

	userID, ok := value.(domain.UserID)
	if !ok {
		return "", xerrors.New("userID not found in context")
	}

	return userID, nil
}
