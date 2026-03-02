package services

import (
	"context"
	"time"

	appContracts "github.com/victorotene80/authentication_api/internal/application/contracts"
	"github.com/victorotene80/authentication_api/internal/application/dto"
	"github.com/victorotene80/authentication_api/internal/domain/aggregates"
	domainContracts "github.com/victorotene80/authentication_api/internal/domain/contracts"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
	"github.com/victorotene80/authentication_api/internal/domain/services/policy"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/infrastructure/persistence/cache"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

var _ appContracts.SessionService = (*sessionService)(nil)

type sessionService struct {
	sessionRepo    repository.SessionRepository
	tokenGen       domainContracts.TokenGenerator
	uow            appContracts.UnitOfWork
	hasher         *utils.SessionKeyHasher
	policy         policy.SessionPolicy
	clock          func() time.Time
	eventPublisher appContracts.EventPublisher
	geoIP          appContracts.GeoIPService
	sessionCache   appContracts.Cache[string, cache.CachedSession]
	auditLogger    appContracts.AuditLogger
}

func NewSessionService(
	sessionRepo repository.SessionRepository,
	tokenGen domainContracts.TokenGenerator,
	uow appContracts.UnitOfWork,
	hasher *utils.SessionKeyHasher,
	policy policy.SessionPolicy,
	clock func() time.Time,
	eventPublisher appContracts.EventPublisher,
	geoIP appContracts.GeoIPService,
	sessionCache appContracts.Cache[string, cache.CachedSession],
	auditLogger appContracts.AuditLogger,
) appContracts.SessionService {
	if clock == nil {
		clock = func() time.Time { return time.Now().UTC() }
	}

	return &sessionService{
		sessionRepo:    sessionRepo,
		tokenGen:       tokenGen,
		uow:            uow,
		hasher:         hasher,
		policy:         policy,
		clock:          clock,
		eventPublisher: eventPublisher,
		geoIP:          geoIP,
		sessionCache:   sessionCache,
		auditLogger:    auditLogger,
	}
}

func (s *sessionService) Create(
	ctx context.Context,
	userID string,
	ipAddress string,
	userAgent string,
	deviceID string,
	deviceFingerprint string,
	deviceName string,
) (appContracts.SessionResult, error) {

	now := s.clock()
	sessionExpiresAt := s.policy.ComputeExpiresAt(now)

	var countryCode, city string
	if s.geoIP != nil && ipAddress != "" {
		if cc, cty, err := s.geoIP.Lookup(ctx, ipAddress); err == nil {
			countryCode, city = cc, cty
		}
	}

	rawSessionKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return appContracts.SessionResult{}, err
	}

	hashedSessionKey := s.hasher.Hash(rawSessionKey)
	tokenHashVO, err := valueobjects.NewSessionTokenHash(hashedSessionKey)
	if err != nil {
		return appContracts.SessionResult{}, err
	}

	var result appContracts.SessionResult

	err = s.uow.WithinTransaction(ctx, func(txCtx context.Context) error {
		session, err := aggregates.NewSession(
			userID,
			tokenHashVO,
			nil,
			ipAddress,
			userAgent,
			deviceFingerprint,
			deviceName,
			countryCode,
			city,
			now,
			sessionExpiresAt,
		)
		if err != nil {
			return err
		}

		accessToken, err := s.tokenGen.GenerateAccess(
			userID,
			session.ID(),
			s.policy.MaxDuration,
		)
		if err != nil {
			return err
		}

		var refreshToken domainContracts.Token
		if s.policy.AllowRefreshToken {
			refreshToken, err = s.tokenGen.GenerateRefresh(
				userID,
				session.ID(),
				s.policy.RefreshTokenDuration,
			)
			if err != nil {
				return err
			}

			hashedRefresh := s.hasher.Hash(refreshToken.Value)
			refreshHashVO, err := valueobjects.NewSessionTokenHash(hashedRefresh)
			if err != nil {
				return err
			}
			session.SetRefreshTokenHash(refreshHashVO)
		}

		if err := s.sessionRepo.Save(txCtx, session); err != nil {
			return err
		}

		if err := s.eventPublisher.Publish(
			txCtx,
			session.PullEvents(),
			map[string]string{
				"aggregate":       "session",
				"action":          "created",
				"ip":              ipAddress,
				"user_agent":      userAgent,
				"device_id":       deviceID,
				"fingerprint":     deviceFingerprint,
				"device_name":     deviceName,
				"country_code":    countryCode,
				"city":            city,
				"session_expires": sessionExpiresAt.UTC().Format(time.RFC3339),
			},
		); err != nil {
			return err
		}

		session.ClearEvents()

		result = appContracts.SessionResult{
			SessionID:    session.ID(),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresAt:    accessToken.ExpiresAt,
		}

		if s.sessionCache != nil {
			cached := cache.MapAggregateToCached(session)
			_ = s.sessionCache.Set(txCtx, session.TokenHash().Value(), &cached)
		}

		return nil
	})

	if err != nil {
		return appContracts.SessionResult{}, err
	}

	if s.auditLogger != nil {
		userIDCopy := userID
		sessionIDCopy := result.SessionID
		ipCopy := ipAddress
		uaCopy := userAgent
		countryCopy := countryCode

		meta := map[string]any{
			"device_id":       deviceID,
			"device_name":     deviceName,
			"fingerprint":     deviceFingerprint,
			"city":            city,
			"session_expires": result.ExpiresAt.UTC().Format(time.RFC3339),
		}

		rec := dto.AuditRecord{
			Action: dto.AuditActionLoginSuccess, 
			UserID: &userIDCopy,
			ActorID:   &userIDCopy,
			SessionID: &sessionIDCopy,

			IPAddress:   &ipCopy,
			UserAgent:   &uaCopy,
			CountryCode: &countryCopy,
			TargetResource: nil, 
			TargetID:       &sessionIDCopy,
			Metadata: meta,
			Success:  true,
		}

		_ = s.auditLogger.Log(ctx, rec) 
	}

	return result, nil
}
