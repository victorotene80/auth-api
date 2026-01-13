package valueobjects

import (
	"strings"
		"github.com/victorotene80/authentication_api/internal/domain"

)

type Role string

const (
	RoleUser      Role = "USER"
	RoleAdmin     Role = "ADMIN"
	RoleModerator Role = "MODERATOR"
)

var roleHierarchy = map[Role][]Role{
	RoleAdmin:     {RoleAdmin, RoleModerator, RoleUser},
	RoleModerator: {RoleModerator, RoleUser},
	RoleUser:      {RoleUser},
}

func NewRole(role string) (Role, error) {
	r := Role(strings.ToUpper(strings.TrimSpace(role)))
	if !r.IsValid() {
		return "", domain.ErrInvalidRole
	}
	return r, nil
}

func (r Role) DefaultUserRole() Role{
	return RoleUser
}

func (r Role) IsValid() bool {
	switch r {
	case RoleUser, RoleAdmin, RoleModerator:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}

func (r Role) HasPermission(required Role) bool {
	allowedRoles, ok := roleHierarchy[r]
	if !ok {
		return false
	}
	for _, role := range allowedRoles {
		if role == required {
			return true
		}
	}
	return false
}
