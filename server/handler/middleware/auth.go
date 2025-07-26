package middleware

import (
	"kakeibo-web-server/domain"
	"kakeibo-web-server/lib/cognito"
	"kakeibo-web-server/lib/ctxdef"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

type CognitoAuthMiddleware struct {
	next      http.Handler
	validator *cognito.Validator
}

func newCognitoAuthMiddleware(next http.Handler, validator *cognito.Validator) http.Handler {
	return &CognitoAuthMiddleware{
		next:      next,
		validator: validator,
	}
}

func (m *CognitoAuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rawAuthorization := r.Header.Get("Authorization")

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

func MakeCognitoAuth(validator *cognito.Validator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return newCognitoAuthMiddleware(next, validator)
	}
}
