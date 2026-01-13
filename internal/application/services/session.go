package services

import (
	"context"
	"time"

	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	"github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

type sessionService struct {
	sessionRepo    repository.SessionRepository
	//sessionCache   appContracts.SessionCache
	tokenGen       contracts.TokenGenerator
	uow            appContracts.UnitOfWork
	hasher         *utils.SessionKeyHasher
	policy         policy.SessionPolicy
	clock          func() time.Time
	eventPublisher appContracts.EventPublisher
}

func NewSessionService(
	sessionRepo repository.SessionRepository,
	//sessionCache appContracts.SessionCache,
	tokenGen contracts.TokenGenerator,
	uow appContracts.UnitOfWork,
	hasher *utils.SessionKeyHasher,
	policy policy.SessionPolicy,
	clock func() time.Time,
	eventPublisher appContracts.EventPublisher,
) contracts.SessionService {
	return &sessionService{
		sessionRepo:    sessionRepo,
		tokenGen:       tokenGen,
		uow:            uow,
		hasher:         hasher,
		policy:         policy,
		clock:          clock,
		eventPublisher: eventPublisher,
	}
}

func (s *sessionService) Create(
	ctx context.Context,
	userID string,
	role valueobjects.Role,
	ipAddress string,
	userAgent string,
	deviceID string,
) (contracts.SessionResult, error) {

	now := s.clock()
	rawRefresh, err := utils.GenerateRandomString(32)
	if err != nil {
		return contracts.SessionResult{}, err
	}

	refreshTokenHash, err := valueobjects.NewSessionTokenHash(s.hasher.Hash(rawRefresh))
	if err != nil {
		return contracts.SessionResult{}, err
	}

	var result contracts.SessionResult

	err = s.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		session, err := aggregates.NewSession(
			userID,
			role,
			refreshTokenHash,
			ipAddress,
			userAgent,
			deviceID,
			now,
			now.Add(s.policy.MaxDuration),
		)

		if err != nil {
			return err
		}

		if err := s.sessionRepo.Save(txCtx, session); err != nil {
			return err
		}

		/*if s.sessionCache != nil {
			_ = s.sessionCache.Set(txCtx, session)
		}*/

		accessToken, err := s.tokenGen.GenerateAccess(
			userID,
			session.ID(),
			role.String(),
			s.policy.MaxDuration,
		)

		if err != nil {
			return err
		}

		refreshToken, err := s.tokenGen.GenerateRefresh(
			userID,
			session.ID(),
			s.policy.RefreshTokenDuration,
		)
		if err != nil {
			return err
		}

		if err := s.eventPublisher.Publish(
			txCtx,
			session.PullEvents(),
			map[string]string{
				"aggregate":  "session",
				"action":     "created",
				"ip":         ipAddress,
				"user_agent": userAgent,
				"device":     deviceID,
			},
		); err != nil {
			return err
		}

		result = contracts.SessionResult{
			SessionID:    session.ID(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    accessToken.ExpiresAt,
		}

		return nil
	})

	return result, err
}
