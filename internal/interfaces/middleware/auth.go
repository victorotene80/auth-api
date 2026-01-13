package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/cache"
	"github.com/victorotene80/authentication_api/internal/infrastructure/services"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

// Context keys
type contextKey string

const (
	ContextUserID contextKey = "userID"
	ContextRole   contextKey = "role"
)

func AuthMiddleware(
	next http.Handler,
	jwtService services.JWTGenerator,
	sessionRepo repository.SessionRepository,
	sessionCache cache.RedisSessionCache,
	userRepo repository.UserRepository,
	hasher *utils.SessionKeyHasher,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// 1️⃣ Extract token
		tokenStr := extractToken(r)
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		// 2️⃣ Validate JWT
		userID, roleFromToken, err := jwtService.Validate(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// 3️⃣ Hash JWT
		hashedToken := valueobjects.NewSessionTokenHash(tokenStr)

		// 4️⃣ Load session from cache
		session, err := sessionCache.Get(ctx, hashedToken)
		if err != nil && !errors.Is(err, infrastructure.ErrSessionNotFound) {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		now := time.Now()

		if session == nil {
			// Cache miss → DB
			session, err = sessionRepo.FindByTokenHash(ctx, hashedToken)
			if err != nil || !session.IsValid(now) {
				http.Error(w, "session invalid", http.StatusUnauthorized)
				return
			}

			// Cache for next request
			_ = sessionCache.Set(ctx, session)
		} else if !session.IsValid(now) {
			// Expired in cache
			_ = sessionCache.Delete(ctx, hashedToken)
			http.Error(w, "session expired", http.StatusUnauthorized)
			return
		}

		// 5️⃣ Optional sensitive endpoint role check
		if isSensitiveEndpoint(r.URL.Path) {
			user, err := userRepo.FindByID(ctx, userID)
			if err != nil || !roleFromTokenHasPermission(roleFromToken, user.User.Role()) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}

		// 6️⃣ Touch session (rolling lastSeen)
		if now.Sub(session.LastSeenAt()) > time.Minute {
			session.Touch(now)
			_ = sessionCache.Set(ctx, session)
		}

		// 7️⃣ Inject into context
		ctx = context.WithValue(ctx, ContextUserID, userID)
		ctx = context.WithValue(ctx, ContextRole, roleFromToken)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helpers

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}

func isSensitiveEndpoint(path string) bool {
	switch path {
	case "/admin", "/withdraw":
		return true
	default:
		return false
	}
}

func roleFromTokenHasPermission(tokenRole string, requiredRole valueobjects.Role) bool {
	r := valueobjects.Role(tokenRole)
	return r.HasPermission(requiredRole)
}

// Typed getters
func UserIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ContextUserID).(string); ok {
		return v
	}
	return ""
}

func RoleFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ContextRole).(string); ok {
		return v
	}
	return ""
}
