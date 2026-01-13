package aggregates

import (
	"fmt"

	"time"

	"github.com/google/uuid"
	"github.com/victorotene80/authentication_api/internal/domain/entities"
	events "github.com/victorotene80/authentication_api/internal/domain/events/types"
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
	role valueobjects.Role,
) (*UserAggregate, error) {

	id := uuid.New().String()

	user := entities.NewUser(
		id,
		email,
		password,
		role,
		firstName,
		lastName,
		middleName,
	)

	agg := &UserAggregate{
		AggregateRoot: NewAggregateRoot(id),
		User:          user,
	}

	event := events.NewUserCreatedEvent(
		id,
		email.String(),
		role.String(),
		firstName,
		lastName,
	)

	agg.RaiseEvent(event)

	return agg, nil
}

func RehydrateUser(
	id, email, hashedPassword, role, firstName, lastName, middleName string,
	failedAttempts int,
	lastFailedAt, lockedAt *time.Time,
	createdAt, updatedAt time.Time,
	version int,
) (*UserAggregate, error) {

	emailVO, err := valueobjects.NewEmail(email)
	if err != nil {
		return nil, err
	}
	passVO, err := valueobjects.NewHashedPassword(hashedPassword)
	if err != nil {
		return nil, err
	}
	roleVO, err := valueobjects.NewRole(role)
	if err != nil {
		return nil, err
	}

	user := entities.NewUser(id, emailVO, passVO, roleVO, firstName, lastName, middleName)
	user.FailedLoginAttempts = failedAttempts
	user.LastFailedAt = lastFailedAt
	user.LockedAt = lockedAt

	agg := &UserAggregate{
		AggregateRoot: &AggregateRoot{
			id:        id,
			version:   version,
			createdAt: createdAt,
			updatedAt: updatedAt,
		},
		User: user,
	}

	return agg, nil
}

func (u *UserAggregate) RecordFailedLogin(now time.Time) {
	u.User.FailedLoginAttempts++
	u.User.LastFailedAt = &now
	if u.User.FailedLoginAttempts >= 5 {
		u.User.LockedAt = &now
	}
}

func (u *UserAggregate) ResetFailedLogins() {
	u.User.FailedLoginAttempts = 0
	u.User.LastFailedAt = nil
}

func (u *UserAggregate) RecordLogin(now time.Time) {
	u.User.SetLastLogin(now)

	event := events.NewUserLogInEvent(
		u.User.ID(),
		u.User.Email().String(),
		u.User.Role().String(),
		u.User.FirstName(),
		u.User.LastName(),
		now,
	)

	u.RaiseEvent(event)
}

func (u *UserAggregate) UnlockUser() {
	u.User.LockedAt = nil
}

func (u *UserAggregate) ChangePassword(
	oldHashed string,
	newHashed valueobjects.Password,
) error {

	if u.User.Password().Value() != oldHashed {
		return fmt.Errorf("old password does not match")
	}

	u.User.UpdatePassword(newHashed)

	u.RaiseEvent(events.NewUserPasswordChangedEvent(
		u.id,
		u.User.Email().String(),
	))

	return nil
}

func (u *UserAggregate) IncrementFailedLogin(now time.Time) {
	u.User.FailedLoginAttempts++
	u.User.LastFailedAt = &now
}

/* TODO
UpdateProfile
ActivateUser
Deactivate
VerifyEmail
*/
