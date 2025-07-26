package middleware

import (
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"kakeibo-web-server/repository"
	"net/http"
)

type DebugAuthMiddleware struct {
	next http.Handler
	repo *repository.Repository
}

func newDebugAuthMiddleware(next http.Handler, repo *repository.Repository) http.Handler {
	return &DebugAuthMiddleware{
		next: next,
		repo: repo,
	}
}

func (m *DebugAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawUserID := r.Header.Get("Debug-User-ID")
	if rawUserID == "" {
		m.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	userID := domain.UserID(rawUserID)

	user, err := m.repo.User.GetByID(ctx, userID)
	if err != nil {
		if err == domain.ErrEntityNotFound {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ctx = ctxdef.WithUserID(ctx, user.ID)

	m.next.ServeHTTP(w, r.WithContext(ctx))
}

func MakeDebugAuth(repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return newDebugAuthMiddleware(next, repo)
	}
}
