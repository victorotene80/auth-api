package aggregates

import (
	"time"

	"github.com/google/uuid"

	events "github.com/victorotene80/authentication_api/internal/domain/events/types"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type SessionStatus string

const (
	SessionActive  SessionStatus = "ACTIVE"
	SessionRevoked SessionStatus = "REVOKED"
	SessionExpired SessionStatus = "EXPIRED"
)

type SessionAggregate struct {
	*AggregateRoot

	id     string
	userID string

	tokenHash        valueobjects.SessionTokenHash
	refreshTokenHash *valueobjects.SessionTokenHash

	ipAddress         string
	deviceFingerprint string
	deviceName        string
	userAgent         string
	countryCode       string
	city              string
	isMFAVerified     bool
	impersonatedBy    *string

	createdAt    time.Time
	lastActiveAt time.Time
	expiresAt    time.Time
	revokedAt    *time.Time
	revokeReason *string
}

// Rehydrate from DB row
func RehydrateSession(
	id, userID string,
	tokenHash string,
	refreshTokenHash *string,
	ip, userAgent, deviceFingerprint, deviceName, countryCode, city string,
	isMFAVerified bool,
	impersonatedBy *string,
	createdAt, lastActiveAt, expiresAt time.Time,
	revokedAt *time.Time,
	revokeReason *string,
	version int,
) (*SessionAggregate, error) {

	currentHash, err := valueobjects.NewSessionTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}

	var refreshHashVO *valueobjects.SessionTokenHash
	if refreshTokenHash != nil {
		h, err := valueobjects.NewSessionTokenHash(*refreshTokenHash)
		if err != nil {
			return nil, err
		}
		refreshHashVO = &h
	}

	ar := NewAggregateRoot(id, version)

	return &SessionAggregate{
		AggregateRoot:     ar,
		id:                id,
		userID:            userID,
		tokenHash:         currentHash,
		refreshTokenHash:  refreshHashVO,
		ipAddress:         ip,
		deviceFingerprint: deviceFingerprint,
		deviceName:        deviceName,
		userAgent:         userAgent,
		countryCode:       countryCode,
		city:              city,
		isMFAVerified:     isMFAVerified,
		impersonatedBy:    impersonatedBy,
		createdAt:         createdAt,
		lastActiveAt:      lastActiveAt,
		expiresAt:         expiresAt,
		revokedAt:         revokedAt,
		revokeReason:      revokeReason,
	}, nil
}

func NewSession(
	userID string,
	tokenHash valueobjects.SessionTokenHash,
	refreshTokenHash *valueobjects.SessionTokenHash,
	ipAddress string,
	userAgent string,
	deviceFingerprint string,
	deviceName string,
	countryCode string,
	city string,
	createdAt time.Time,
	expiresAt time.Time,
) (*SessionAggregate, error) {

	id := uuid.NewString()
	ar := NewAggregateRoot(id, 0)

	s := &SessionAggregate{
		AggregateRoot:     ar,
		id:                id,
		userID:            userID,
		tokenHash:         tokenHash,
		refreshTokenHash:  refreshTokenHash,
		ipAddress:         ipAddress,
		deviceFingerprint: deviceFingerprint,
		deviceName:        deviceName,
		userAgent:         userAgent,
		countryCode:       countryCode,
		city:              city,
		isMFAVerified:     false,
		impersonatedBy:    nil,
		createdAt:         createdAt,
		lastActiveAt:      createdAt,
		expiresAt:         expiresAt,
		revokedAt:         nil,
		revokeReason:      nil,
	}

	s.RaiseEvent(events.NewSessionCreatedEvent(
		s.id,
		s.userID,
		ipAddress,
		userAgent,
		deviceName,
	))

	return s, nil
}

func (s *SessionAggregate) Status(now time.Time) SessionStatus {
	if s.revokedAt != nil {
		return SessionRevoked
	}
	if now.After(s.expiresAt) {
		return SessionExpired
	}
	return SessionActive
}

func (s *SessionAggregate) IsValid(now time.Time) bool {
	return s.Status(now) == SessionActive
}

func (s *SessionAggregate) Touch(now time.Time) {
	if s.Status(now) == SessionActive {
		s.lastActiveAt = now
	}
}

func (s *SessionAggregate) Revoke(now time.Time, reason string) {
	if s.revokedAt != nil {
		return
	}
	s.revokedAt = &now
	s.revokeReason = &reason
	s.RaiseEvent(events.NewSessionRevokedEvent(s.id, s.userID, reason))
}

// Simplified rotate: just update token hash and bump activity.
// We no longer persist previous hash / rotation id because DB doesn’t have them.
func (s *SessionAggregate) RotateKey(
	newTokenHash valueobjects.SessionTokenHash,
	now time.Time,
) {
	if s.Status(now) != SessionActive {
		return
	}

	s.tokenHash = newTokenHash
	s.lastActiveAt = now

	s.RaiseEvent(events.NewSessionAccessedEvent(s.id, s.userID))
}

func (s *SessionAggregate) SetRefreshTokenHash(hash valueobjects.SessionTokenHash) {
	s.refreshTokenHash = &hash
}

// Getters
func (s *SessionAggregate) ID() string                          { return s.id }
func (s *SessionAggregate) UserID() string                      { return s.userID }
func (s *SessionAggregate) TokenHash() valueobjects.SessionTokenHash {
	return s.tokenHash
}
func (s *SessionAggregate) RefreshTokenHash() *valueobjects.SessionTokenHash {
	return s.refreshTokenHash
}
func (s *SessionAggregate) CreatedAt() time.Time     { return s.createdAt }
func (s *SessionAggregate) LastActiveAt() time.Time  { return s.lastActiveAt }
func (s *SessionAggregate) ExpiresAt() time.Time     { return s.expiresAt }
func (s *SessionAggregate) RevokedAt() *time.Time    { return s.revokedAt }
func (s *SessionAggregate) RevokeReason() *string    { return s.revokeReason }
func (s *SessionAggregate) IPAddress() string        { return s.ipAddress }
func (s *SessionAggregate) UserAgent() string        { return s.userAgent }
func (s *SessionAggregate) DeviceFingerprint() string { return s.deviceFingerprint }
func (s *SessionAggregate) DeviceName() string       { return s.deviceName }
func (s *SessionAggregate) CountryCode() string      { return s.countryCode }
func (s *SessionAggregate) City() string             { return s.city }
func (s *SessionAggregate) IsMFAVerified() bool      { return s.isMFAVerified }
func (s *SessionAggregate) ImpersonatedBy() *string  { return s.impersonatedBy }
