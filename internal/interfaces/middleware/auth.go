package middleware

import (
	"context"
	"net/http"

	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
)

type ctxKey string

const (
	ctxKeyUserID    ctxKey = "auth_user_id"
	ctxKeySessionID ctxKey = "auth_session_id"
)

type AuthMiddleware struct {
	authService appContracts.AuthService
}

func NewAuthMiddleware(authService appContracts.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := extractBearerToken(r)
		if tokenStr == "" {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		authCtx, err := m.authService.Authenticate(r.Context(), tokenStr)
		if err != nil {
			// Map app error → HTTP status (401/403)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUserID, authCtx.UserID)
		ctx = context.WithValue(ctx, ctxKeySessionID, authCtx.SessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractBearerToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

func UserIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyUserID).(string); ok {
		return v
	}
	return ""
}

func SessionIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeySessionID).(string); ok {
		return v
	}
	return ""
}