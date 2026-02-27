package aggregates

import (
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"

	"github.com/victorotene80/authentication_api/internal/domain/entities"
	"github.com/victorotene80/authentication_api/internal/domain/events/types"
	"github.com/victorotene80/authentication_api/internal/domain/services"
	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
)

type UserAggregate struct {
	*AggregateRoot
	User *entities.User
}


func NewUserAggregate(
	email valueobjects.Email,
	password valueobjects.Password,
	firstName, lastName, middleName string,
	now time.Time,
) (*UserAggregate, error) {

	id := uuid.NewString()

	user := entities.NewUserForRegistration(
		id,
		email,
		password,
		firstName,
		lastName,
		middleName,
		now,
	)

	agg := &UserAggregate{
		AggregateRoot: NewAggregateRoot(id, 0),
		User:          user,
	}

	agg.RaiseEvent(types.NewUserCreatedEvent(
		id,
		email.String(),
		user.Status().String(),
		firstName,
		lastName,
	))

	return agg, nil
}


func NewUserAggregateFromDBRow(
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
	version int,
) (*UserAggregate, error) {

	user := entities.NewUserFromDB(
		id,
		email,
		password,
		status,
		firstName,
		lastName,
		middleName,
		emailVerified,
		emailVerifiedAt,
		passwordChangedAt,
		passwordExpiresAt,
		lockedUntil,
		lastLoginAt,
		lastActiveAt,
		deletedAt,
		lastLoginIP,
		failedAttempts,
		createdAt,
		updatedAt,
	)

	return &UserAggregate{
		AggregateRoot: NewAggregateRoot(id, version),
		User:          user,
	}, nil
}

// RehydrateUser is a convenience that takes raw strings from DB and builds VOs.
func RehydrateUser(
	id, emailStr, hashedPassword, statusStr,
	firstName, lastName, middleName string,
	emailVerified bool,
	emailVerifiedAt, passwordChangedAt, passwordExpiresAt, lockedUntil,
	lastLoginAt, lastActiveAt, deletedAt *time.Time,
	lastLoginIP net.IP,
	failedAttempts int,
	createdAt, updatedAt time.Time,
	version int,
) (*UserAggregate, error) {

	emailVO, err := valueobjects.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	passVO, err := valueobjects.NewHashedPassword(hashedPassword)
	if err != nil {
		return nil, err
	}

	status := valueobjects.UserStatus(statusStr)

	return NewUserAggregateFromDBRow(
		id,
		emailVO,
		passVO,
		status,
		firstName,
		lastName,
		middleName,
		emailVerified,
		emailVerifiedAt,
		passwordChangedAt,
		passwordExpiresAt,
		lockedUntil,
		lastLoginAt,
		lastActiveAt,
		deletedAt,
		lastLoginIP,
		failedAttempts,
		createdAt,
		updatedAt,
		version,
	)
}

func (u *UserAggregate) RecordLogin(now time.Time, ip net.IP) {
	u.User.RecordLogin(now, ip)

	u.RaiseEvent(types.NewUserLogInEvent(
		u.User.ID(),
		u.User.Email().String(),
		now,
	))
}

func (u *UserAggregate) ChangePassword(
	currentHashed string,
	newHashed valueobjects.Password,
	now time.Time,
) error {
	if u.User.Password().Value() != currentHashed {
		return fmt.Errorf("old password does not match")
	}

	u.User.ChangePassword(newHashed, now, false)

	u.RaiseEvent(types.NewUserPasswordChangedEvent(
		u.ID(),
		u.User.Email().String(),
	))

	return nil
}

func (u *UserAggregate) RecordFailedLogin(now time.Time, lockSvc *services.AccountLockService) {
	u.User.IncrementFailedLogin()

	if lockSvc.ShouldLock(u.User.FailedLoginAttempts()) {
		lockedUntil := lockSvc.ComputeLockedUntil(now)
		u.User.LockUntil(lockedUntil)

		u.RaiseEvent(types.NewUserLockedEvent(
			u.User.ID(),
			u.User.Email().String(),
			lockedUntil,
		))
	}

	u.RaiseEvent(types.NewUserLoginFailedEvent(
		u.User.ID(),
		u.User.Email().String(),
		now,
	))
}

func (u *UserAggregate) EnsureNotLocked(now time.Time, lockSvc *services.AccountLockService) bool {
	if lockSvc.IsLocked(u.User.LockedUntil(), now) {
		return false
	}

	u.User.UnlockIfExpired(now)
	return true
}