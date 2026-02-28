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

type contextKey string

const (
	ContextUserID    contextKey = "userID"
	ContextSessionID contextKey = "sessionID"
)

func AuthMiddleware(
    next http.Handler,
    jwtService *services.JWTGenerator,
    sessionRepo repository.SessionRepository,
    sessionCache cache.RedisSessionCache,
    hasher *utils.SessionKeyHasher,
) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()

        tokenStr := extractToken(r)
        if tokenStr == "" {
            http.Error(w, "missing token", http.StatusUnauthorized)
            return
        }

        userID, sessionID, err := jwtService.ValidateAccess(tokenStr)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        hashed := hasher.Hash(tokenStr)
        tokenHashVO, err := valueobjects.NewSessionTokenHash(hashed)
        if err != nil {
            http.Error(w, "internal error", http.StatusInternalServerError)
            return
        }

        now := time.Now().UTC()

        session, err := sessionCache.Get(ctx, tokenHashVO.Value())
        if err != nil && !errors.Is(err, infrastructure.ErrSessionNotFound) {
            http.Error(w, "internal error", http.StatusInternalServerError)
            return
        }

        if session == nil {
            session, err = sessionRepo.FindByTokenHash(ctx, tokenHashVO, now)
            if err != nil || !session.IsValid(now) {
                http.Error(w, "session invalid", http.StatusUnauthorized)
                return
            }

        } else if !session.IsValid(now) {
            _ = sessionCache.Delete(ctx, tokenHashVO.Value())
            http.Error(w, "session expired", http.StatusUnauthorized)
            return
        }

        if now.Sub(session.LastActiveAt()) > time.Minute {
            session.Touch(now)
        }

        ctx = context.WithValue(ctx, ContextUserID, userID)
        ctx = context.WithValue(ctx, ContextSessionID, sessionID)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if len(auth) > 7 && auth[:7] == "Bearer " {
		return auth[7:]
	}
	return ""
}
