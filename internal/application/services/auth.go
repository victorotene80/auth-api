// internal/application/services/auth_service.go
package services

import (
	"context"
	"time"

	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	appErrors "github.com/victorotene80/authentication_api/internal/application"

	domainContracts "github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
)

var _ appContracts.AuthService = (*AuthServiceImpl)(nil)

type AuthServiceImpl struct {
	tokenGen      domainContracts.TokenGenerator
	sessionRepo   repository.SessionRepository
	sessionPolicy policy.SessionPolicy
}

func NewAuthService(
	tokenGen domainContracts.TokenGenerator,
	sessionRepo repository.SessionRepository,
	sessionPolicy policy.SessionPolicy,
) *AuthServiceImpl {
	return &AuthServiceImpl{
		tokenGen:      tokenGen,
		sessionRepo:   sessionRepo,
		sessionPolicy: sessionPolicy,
	}
}

func (s *AuthServiceImpl) Authenticate(
	ctx context.Context,
	accessToken string,
) (appContracts.AuthContext, error) {

	userID, sessionID, err := s.tokenGen.ValidateAccess(accessToken)
	if err != nil {
		return appContracts.AuthContext{}, appErrors.ErrUnauthenticated
	}

	now := time.Now().UTC()

	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return appContracts.AuthContext{}, appErrors.ErrUnauthenticated
	}

	if s.sessionPolicy.IsSessionExpired(session.CreatedAt(), session.LastActiveAt(), now) {
		session.Revoke(now, "expired_by_policy")
		_ = s.sessionRepo.Save(ctx, session)

		return appContracts.AuthContext{}, appErrors.ErrUnauthenticated
	}

	if !session.IsValid(now) {
		return appContracts.AuthContext{}, appErrors.ErrUnauthenticated
	}

	session.Touch(now)
	_ = s.sessionRepo.Save(ctx, session)

	return appContracts.AuthContext{
		UserID:    userID,
		SessionID: sessionID,
	}, nil
}