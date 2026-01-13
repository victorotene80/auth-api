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
	id                string
	userID            string
	tokenHash         valueobjects.SessionTokenHash
	previousTokenHash *valueobjects.SessionTokenHash
	rotationID        *string
	ipAddress         string
	deviceID          string
	userAgent         string
	role              valueobjects.Role 
	status            SessionStatus
	createdAt         time.Time
	lastSeenAt        time.Time
	expiresAt         time.Time
	revokedAt         *time.Time
}

func NewSession(
	userID string,
	role valueobjects.Role,
	tokenHash valueobjects.SessionTokenHash,
	ipAddress, userAgent, deviceID string,
	now, expiresAt time.Time,
) (*SessionAggregate, error) {

	id := uuid.NewString()

	s := &SessionAggregate{
		AggregateRoot: NewAggregateRoot(id),
		id:            id,
		userID:        userID,
		role:          role, 
		tokenHash:     tokenHash,
		ipAddress:     ipAddress,
		userAgent:     userAgent,
		deviceID:      deviceID,
		status:        SessionActive,
		createdAt:     now,
		lastSeenAt:    now,
		expiresAt:     expiresAt,
	}

	event := events.NewSessionCreatedEvent(id, userID, ipAddress, userAgent, expiresAt.Format(time.RFC3339))
	s.RaiseEvent(event)

	return s, nil
}

func RehydrateSession(
	id, userID string,
	role valueobjects.Role, 
	tokenHash string,
	prevTokenHash, rotationID *string,
	ip, agent, deviceID, status string,
	createdAt, lastSeenAt, expiresAt time.Time,
	revokedAt *time.Time,
	version int,
) (*SessionAggregate, error) {

	currentHash, err := valueobjects.NewSessionTokenHash(tokenHash)
	if err != nil {
		return nil, err
	}

	var prevHashVO *valueobjects.SessionTokenHash
	if prevTokenHash != nil {
		h, err := valueobjects.NewSessionTokenHash(*prevTokenHash)
		if err != nil {
			return nil, err
		}
		prevHashVO = &h
	}

	ar := NewAggregateRoot(id)
	ar.version = version

	return &SessionAggregate{
		AggregateRoot:     ar,
		id:                id,
		userID:            userID,
		role:              role,
		tokenHash:         currentHash,
		previousTokenHash: prevHashVO,
		rotationID:        rotationID,
		ipAddress:         ip,
		deviceID:          deviceID,
		userAgent:         agent,
		status:            SessionStatus(status),
		createdAt:         createdAt,
		lastSeenAt:        lastSeenAt,
		expiresAt:         expiresAt,
		revokedAt:         revokedAt,
	}, nil
}

func (s *SessionAggregate) IsValid(now time.Time) bool {
	return s.status == SessionActive && now.Before(s.expiresAt)
}

func (s *SessionAggregate) EvaluateExpiry(now time.Time) bool {
	if s.status == SessionActive && now.After(s.expiresAt) {
		s.status = SessionExpired
		s.RaiseEvent(events.NewSessionExpiredEvent(s.id, s.userID))
		return true
	}
	return false
}

func (s *SessionAggregate) Touch(now time.Time) {
	if s.status == SessionActive {
		s.lastSeenAt = now
	}
}

func (s *SessionAggregate) Revoke(now time.Time, reason string) error {
	if s.status == SessionRevoked {
		return nil
	}

	s.status = SessionRevoked
	s.revokedAt = &now
	s.RaiseEvent(events.NewSessionRevokedEvent(s.id, s.userID, reason))
	return nil
}

func (s *SessionAggregate) RotateKey(
	newTokenHash valueobjects.SessionTokenHash,
	rotationID string,
	now time.Time,
) error {
	if s.status != SessionActive {
		return nil
	}
	if s.rotationID != nil && *s.rotationID == rotationID {
		return nil
	}

	old := s.tokenHash
	s.previousTokenHash = &old
	s.tokenHash = newTokenHash
	s.rotationID = &rotationID
	s.lastSeenAt = now

	s.RaiseEvent(events.NewSessionAccessedEvent(s.id, s.userID))

	return nil
}

func (s *SessionAggregate) IsRevoked() bool {
	return s.status == SessionRevoked
}

func (s *SessionAggregate) ID() string                               { return s.id }
func (s *SessionAggregate) UserID() string                           { return s.userID }
func (s *SessionAggregate) TokenHash() valueobjects.SessionTokenHash { return s.tokenHash }
func (s *SessionAggregate) PreviousTokenHash() *valueobjects.SessionTokenHash {
	return s.previousTokenHash
}
func (s *SessionAggregate) RotationID() *string   { return s.rotationID }
func (s *SessionAggregate) Status() SessionStatus { return s.status }
func (s *SessionAggregate) CreatedAt() time.Time  { return s.createdAt }
func (s *SessionAggregate) LastSeenAt() time.Time { return s.lastSeenAt }
func (s *SessionAggregate) ExpiresAt() time.Time  { return s.expiresAt }
func (s *SessionAggregate) RevokedAt() *time.Time { return s.revokedAt }
func (s *SessionAggregate) IPAddress() string     { return s.ipAddress }
func (s *SessionAggregate) UserAgent() string     { return s.userAgent }
func (s *SessionAggregate) DeviceID() string      { return s.deviceID }
func (s *SessionAggregate) Version() int          { return s.AggregateRoot.version }
func (s *SessionAggregate) Role() valueobjects.Role {
	return s.role
}
