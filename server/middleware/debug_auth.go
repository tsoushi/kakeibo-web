package middleware

import (
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/ctxdef"
	"net/http"
)

type DebugAuthMiddleware struct {
	next http.Handler
}

func newDebugAuthMiddleware(next http.Handler) http.Handler {
	return &DebugAuthMiddleware{
		next: next,
	}
}

func (m *DebugAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.Header.Get("Debug-User-ID")
	if userID != "" {
		ctx = ctxdef.WithUserID(ctx, domain.UserID(userID))
	}
	m.next.ServeHTTP(w, r.WithContext(ctx))
}

func DebugAuth(next http.Handler) http.Handler {
	return newDebugAuthMiddleware(next)
}
