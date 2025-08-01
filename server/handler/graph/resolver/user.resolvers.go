package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.72

import (
	"context"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/handler/graph"
	"kakeibo-web-server/lib/ctxdef"

	"golang.org/x/xerrors"
)

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context) (*domain.User, error) {
	userID, err := ctxdef.UserID(ctx)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	user, err := r.usecase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return user, nil
}

// ID is the resolver for the ID field.
func (r *userResolver) ID(ctx context.Context, obj *domain.User) (string, error) {
	return string(obj.ID), nil
}

// User returns graph.UserResolver implementation.
func (r *Resolver) User() graph.UserResolver { return &userResolver{r} }

type userResolver struct{ *Resolver }
