package valueobjects

type UserStatus string

const (
	UserStatusPendingVerification UserStatus = "pending_verification"
	UserStatusActive              UserStatus = "active"
	UserStatusSuspended           UserStatus = "suspended"
	UserStatusDeactivated         UserStatus = "deactivated"
	UserStatusLocked              UserStatus = "locked"
)

func (s UserStatus) String() string {
	return string(s)
}

func (s UserStatus) IsLocked() bool {
	return s == UserStatusLocked
}

func(s UserStatus) IsActive() bool {
	return s == UserStatusActive
}

func (s UserStatus) IsSuspended() bool {
	return s == UserStatusSuspended
}

func (s UserStatus) IsDeactivated() bool {
	return s == UserStatusDeactivated
}

func (s UserStatus) IsPendingVerification() bool {
	return s == UserStatusPendingVerification
}

