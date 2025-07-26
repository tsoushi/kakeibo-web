package repository

import (
	"context"
	"errors"
	"kakeibo-web-server/domain"

	"github.com/gocraft/dbr/v2"
	"golang.org/x/xerrors"
)

const usertableName = "user"

type UserRepository struct {
	sess *dbr.Session
}

func NewUserRepository(sess *dbr.Session) *UserRepository {
	return &UserRepository{
		sess: sess,
	}
}

func (r *UserRepository) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	var user domain.User
	runner := getRunner(ctx, r.sess)
	err := runner.Select("*").From(usertableName).Where("id = ?", id).LoadOneContext(ctx, &user)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByName(ctx context.Context, name string) (*domain.User, error) {
	var user domain.User
	runner := getRunner(ctx, r.sess)
	err := runner.Select("*").From(usertableName).Where("name = ?", name).LoadOneContext(ctx, &user)
	if err != nil {
		if errors.Is(err, dbr.ErrNotFound) {
			return nil, domain.ErrEntityNotFound
		}
		return nil, xerrors.Errorf("failed to get user by name: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) Insert(ctx context.Context, user *domain.User) (*domain.User, error) {
	runner := getRunner(ctx, r.sess)
	_, err := runner.InsertInto(usertableName).Columns("id", "name").Record(user).Exec()
	if err != nil {
		return nil, xerrors.Errorf("failed to insert user: %w", err)
	}

	return user, nil
}
