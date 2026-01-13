package entities

import (
	"errors"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/valueobjects"
	"github.com/victorotene80/authentication_api/internal/shared/utils"
)

type User struct {
	id                  string
	email               valueobjects.Email
	password            valueobjects.Password
	firstName           string
	lastName            string
	middleName          string
	role                valueobjects.Role
	isActive            bool
	isVerified          bool
	lastLoginAt         *time.Time
	FailedLoginAttempts int
	LastFailedAt        *time.Time
	LockedAt            *time.Time
	createdAt           time.Time
	updatedAt           time.Time
	deletedAt           *time.Time
}

func NewUser(
	id string,
	email valueobjects.Email,
	password valueobjects.Password,
	role valueobjects.Role,
	firstName, lastName, middleName string,
) *User {
	now := utils.NowUTC()
	return &User{
		id:         id,
		email:      email,
		password:   password,
		firstName:  firstName,
		lastName:   lastName,
		middleName: middleName,
		role:       role,
		isActive:   true,
		isVerified: false,
		createdAt:  now,
		updatedAt:  now,
	}
}

func (u *User) UpdatePassword(password valueobjects.Password) {
	u.password = password
	u.updatedAt = utils.NowUTC()
}

func (u *User) UpdateProfile(firstName, lastName, middleName string) {
	u.firstName = firstName
	u.lastName = lastName
	u.middleName = middleName
	u.updatedAt = utils.NowUTC()
}

func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = utils.NowUTC()
}

func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = utils.NowUTC()
}

func (u *User) VerifyEmail() {
	u.isVerified = true
	u.updatedAt = utils.NowUTC()
}

func (u *User) SetLastLogin(now time.Time) {
	u.lastLoginAt = &now
	u.updatedAt = utils.NowUTC()
}

func (u *User) Delete(now time.Time) error {
	if u.deletedAt != nil {
		return errors.New("user already deleted")
	}

	u.deletedAt = &now
	u.isActive = false
	u.updatedAt = now
	return nil
}

func (u *User) IsDeleted() bool {
	return u.deletedAt != nil
}

func (u *User) DeletedAt() *time.Time {
	return u.deletedAt
}

func (u *User) ID() string                      { return u.id }
func (u *User) Email() valueobjects.Email       { return u.email }
func (u *User) Password() valueobjects.Password { return u.password }
func (u *User) Role() valueobjects.Role         { return u.role }
func (u *User) FirstName() string               { return u.firstName }
func (u *User) LastName() string                { return u.lastName }
func (u *User) MiddleName() string              { return u.middleName }
func (u *User) IsActive() bool                  { return u.isActive }
func (u *User) IsVerified() bool                { return u.isVerified }
func (u *User) LastLoginAt() *time.Time         { return u.lastLoginAt }
func (u *User) CreatedAt() time.Time            { return u.createdAt }
func (u *User) UpdatedAt() time.Time            { return u.updatedAt }