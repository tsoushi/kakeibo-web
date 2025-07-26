package middleware

import (
	"errors"
	"fmt"
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/cognito"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/repository"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

type CognitoAuthMiddleware struct {
	next      http.Handler
	validator *cognito.Validator
	repo      *repository.Repository
}

func newCognitoAuthMiddleware(next http.Handler, validator *cognito.Validator, repo *repository.Repository) http.Handler {
	return &CognitoAuthMiddleware{
		next:      next,
		validator: validator,
		repo:      repo,
	}
}

func (m *CognitoAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawAuthorization := r.Header.Get("Authorization")
	if rawAuthorization == "" {
		m.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	tokenString, err := extractTokenFromAuthorization(rawAuthorization)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	claims, err := m.validator.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID := domain.UserID(claims.Subject)

	_, err = m.repo.User.GetByID(ctx, userID)
	if err != nil && !errors.Is(err, domain.ErrEntityNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, domain.ErrEntityNotFound) {
		_, err = m.repo.User.Insert(ctx, domain.NewUser(userID, claims.Username))
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to insert user: %v", err), http.StatusInternalServerError)
			return
		}
	}

	ctx = ctxdef.WithUserID(ctx, domain.UserID(claims.Subject))

	m.next.ServeHTTP(w, r.WithContext(ctx))
}

func extractTokenFromAuthorization(rawAuthorization string) (string, error) {
	const prefix = "Bearer "
	if !strings.HasPrefix(rawAuthorization, prefix) {
		return "", xerrors.New("invalid authorization header format")
	}

	return strings.TrimPrefix(rawAuthorization, prefix), nil
}

func MakeCognitoAuth(validator *cognito.Validator, repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return newCognitoAuthMiddleware(next, validator, repo)
	}
}
