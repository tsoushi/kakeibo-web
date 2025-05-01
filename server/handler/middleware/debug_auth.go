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

	userName := r.Header.Get("Debug-User-Name")
	if userName == "" {
		m.next.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	password := r.Header.Get("Debug-User-Password")

	user, err := m.repo.User.GetByName(ctx, userName)
	if err != nil {
		if err == domain.ErrEntityNotFound {
			http.Error(w, "Invalid user name or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	ok, err := user.HashedPassword.Compare(password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if !ok {
		http.Error(w, "Invalid user name or password", http.StatusUnauthorized)
		return
	}

	ctx = ctxdef.WithUserID(ctx, domain.UserID(user.ID))

	m.next.ServeHTTP(w, r.WithContext(ctx))
}

func MakeDebugAuth(repo *repository.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return newDebugAuthMiddleware(next, repo)
	}
}
