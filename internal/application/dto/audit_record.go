package dto

type AuditAction string

const (
	AuditActionLoginSuccess        AuditAction = "login_success"
	AuditActionLoginFailed         AuditAction = "login_failed"
	AuditActionLogout              AuditAction = "logout"
	AuditActionSessionRevoked      AuditAction = "session_revoked"
	AuditActionSessionCreated      AuditAction = "session_created"
	AuditActionSessionRotated      AuditAction = "session_rotated"
	AuditActionPasswordChanged     AuditAction = "password_changed"
	AuditActionPasswordResetReq    AuditAction = "password_reset_requested"
	AuditActionPasswordResetDone   AuditAction = "password_reset_completed"
	AuditActionEmailVerified       AuditAction = "email_verified"
	AuditActionSuspiciousActivity  AuditAction = "suspicious_activity"
)

type AuditRecord struct {
	Action AuditAction
	UserID    *string
	ActorID   *string
	APIKeyID  *string
	SessionID *string
	OrganizationID *string
	IPAddress   *string
	UserAgent   *string
	CountryCode *string
	TargetResource *string
	TargetID       *string
	Metadata map[string]any
	Success       bool
	FailureReason *string
}
