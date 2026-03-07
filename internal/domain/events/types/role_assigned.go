package types

import "github.com/victorotene80/authentication_api/internal/domain/events"

type UserRoleAssignedPayload struct {
	UserID   string
	RoleID   string
	RoleSlug string
	Version  int
}

const UserRoleAssignedEventName = "user.role_assigned"

func NewUserRoleAssignedEvent(
	userID string,
	roleID string,
	roleSlug string,
	version int,
) events.DomainEvent {
	return events.NewEvent(
		UserRoleAssignedEventName,
		userID,
		UserRoleAssignedPayload{
			UserID:   userID,
			RoleID:   roleID,
			RoleSlug: roleSlug,
			Version:  version,
		},
		nil,
	)
}