package entities

import (
	"errors"
	"net"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type User struct {
	id                    string
	email                 valueobjects.Email
	password              valueobjects.Password
	passwordChangedAt     *time.Time
	passwordExpiresAt     *time.Time
	requirePasswordChange bool
	firstName             string
	lastName              string
	middleName            string
	status                valueobjects.UserStatus
	emailVerified         bool
	emailVerifiedAt       *time.Time
	failedLoginAttempts   int
	lockedUntil           *time.Time
	lastLoginAt           *time.Time
	lastLoginIP           net.IP
	lastActiveAt          *time.Time
	createdAt             time.Time
	updatedAt             time.Time
	deletedAt             *time.Time
}

func NewUserForRegistration(
	id string,
	email valueobjects.Email,
	password valueobjects.Password,
	firstName, lastName, middleName string,
	now time.Time,
) *User {
	return &User{
		id:                    id,
		email:                 email,
		password:              password,
		firstName:             firstName,
		lastName:              lastName,
		middleName:            middleName,
		status:                valueobjects.UserStatusPendingVerification, // or Active if you prefer
		emailVerified:         false,
		emailVerifiedAt:       nil,
		passwordChangedAt:     nil,
		passwordExpiresAt:     nil,
		requirePasswordChange: false,
		failedLoginAttempts:   0,
		lockedUntil:           nil,
		lastLoginAt:           nil,
		lastLoginIP:           nil,
		lastActiveAt:          nil,
		createdAt:             now,
		updatedAt:             now,
		deletedAt:             nil,
	}
}


func NewUserFromDB(
	id string,
	email valueobjects.Email,
	password valueobjects.Password,
	status valueobjects.UserStatus,
	firstName, lastName, middleName string,
	emailVerified bool,
	emailVerifiedAt, passwordChangedAt, passwordExpiresAt, lockedUntil,
	lastLoginAt, lastActiveAt, deletedAt *time.Time,
	lastLoginIP net.IP,
	failedAttempts int,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:                    id,
		email:                 email,
		password:              password,
		status:                status,
		firstName:             firstName,
		lastName:              lastName,
		middleName:            middleName,
		emailVerified:         emailVerified,
		emailVerifiedAt:       emailVerifiedAt,
		passwordChangedAt:     passwordChangedAt,
		passwordExpiresAt:     passwordExpiresAt,
		requirePasswordChange: false,
		failedLoginAttempts:   failedAttempts,
		lockedUntil:           lockedUntil,
		lastLoginAt:           lastLoginAt,
		lastLoginIP:           lastLoginIP,
		lastActiveAt:          lastActiveAt,
		createdAt:             createdAt,
		updatedAt:             updatedAt,
		deletedAt:             deletedAt,
	}
}

func (u *User) MarkEmailVerified(at time.Time) {
	u.emailVerified = true
	u.emailVerifiedAt = &at
}

func (u *User) UpdateProfile(firstName, lastName, middleName string) {
	u.firstName = firstName
	u.lastName = lastName
	u.middleName = middleName
}

func (u *User) ChangePassword(newPassword valueobjects.Password, changedAt time.Time, requireChange bool) {
	u.password = newPassword
	u.passwordChangedAt = &changedAt
	u.requirePasswordChange = requireChange
}

func (u *User) SetStatus(status valueobjects.UserStatus) {
	u.status = status
}

func (u *User) RecordLogin(at time.Time, ip net.IP) {
	u.lastLoginAt = &at
	u.lastLoginIP = ip
	u.lastActiveAt = &at
	u.failedLoginAttempts = 0
	u.lockedUntil = nil
}

func (u *User) RecordActivity(at time.Time) {
	u.lastActiveAt = &at
}

func (u *User) RecordFailedLogin(at time.Time) {
	u.failedLoginAttempts++
}

func (u *User) ApplyLock(lockedUntil time.Time) {
	u.lockedUntil = &lockedUntil
	u.status = valueobjects.UserStatusLocked
}

func (u *User) ClearLock() {
	u.lockedUntil = nil
	if u.status == valueobjects.UserStatusLocked {
		u.status = valueobjects.UserStatusActive
	}
}

func (u *User) SoftDelete(at time.Time) error {
	if u.deletedAt != nil {
		return errors.New("user already deleted")
	}
	u.deletedAt = &at
	return nil
}

func (u *User) IncrementFailedLogin() error {
	if u.lockedUntil != nil {
		return errors.New("account is currently locked")
	}
	u.failedLoginAttempts++
	return nil
}


func (u *User) LockUntil(until time.Time) {
	u.lockedUntil = &until
	u.status = valueobjects.UserStatusLocked
}

func (u *User) UnlockIfExpired(now time.Time) {
	if u.lockedUntil != nil && !now.Before(*u.lockedUntil) {
		u.lockedUntil = nil
		if u.status == valueobjects.UserStatusLocked {
			u.status = valueobjects.UserStatusActive
		}
		u.failedLoginAttempts = 0 // optional, but recommended after lock expires
	}
}

func (u *User) ID() string                      { return u.id }
func (u *User) Email() valueobjects.Email       { return u.email }
func (u *User) Password() valueobjects.Password { return u.password }
func (u *User) Status() valueobjects.UserStatus { return u.status }
func (u *User) FirstName() string               { return u.firstName }
func (u *User) LastName() string                { return u.lastName }
func (u *User) MiddleName() string              { return u.middleName }
func (u *User) EmailVerified() bool             { return u.emailVerified }
func (u *User) EmailVerifiedAt() *time.Time     { return u.emailVerifiedAt }
func (u *User) FailedLoginAttempts() int        { return u.failedLoginAttempts }
func (u *User) LockedUntil() *time.Time         { return u.lockedUntil }
func (u *User) LastLoginAt() *time.Time         { return u.lastLoginAt }
func (u *User) LastLoginIP() net.IP             { return u.lastLoginIP }
func (u *User) LastActiveAt() *time.Time        { return u.lastActiveAt }
func (u *User) CreatedAt() time.Time            { return u.createdAt }
func (u *User) UpdatedAt() time.Time            { return u.updatedAt }
func (u *User) DeletedAt() *time.Time           { return u.deletedAt }
