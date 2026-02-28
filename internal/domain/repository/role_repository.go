package repository

import (
	"context"
	"time"

	"github.com/victorotene80/authentication_api/internal/domain/entities"
)

type RoleRepository interface {
	FindBySlug(ctx context.Context, slug string) (*entities.Role, error)
	AssignRole(
		ctx context.Context,
		userID string,
		roleID string,
		organizationID *string,
		grantedBy *string,
		expiresAt *time.Time,
	) error
}
