// internal/infrastructure/persistence/role_repository.go
package persistence

import (
	"context"
	"database/sql"
	"time"
	"github.com/victorotene80/authentication_api/internal/domain/entities"
	"github.com/victorotene80/authentication_api/internal/domain/repository"
)

type PgRoleRepository struct {
	db *sql.DB
}
func NewPgRoleRepository(db *sql.DB) repository.RoleRepository {
	return &PgRoleRepository{db: db}
}

func (r *PgRoleRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*entities.Role, error) {

	const q = `
        SELECT
            id,
            name,
            slug,
            description,
            parent_id,
            is_system,
            created_at,
            updated_at
        FROM auth.roles
        WHERE slug = $1
        LIMIT 1;
    `

	row := r.db.QueryRowContext(ctx, q, slug)

	var (
		role     entities.Role
		desc     sql.NullString
		parentID sql.NullString
	)

	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Slug,
		&desc,
		&parentID,
		&role.IsSystem,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // "not found"; caller decides behavior
		}
		return nil, err
	}

	if desc.Valid {
		role.Description = &desc.String
	}
	if parentID.Valid {
		role.ParentID = &parentID.String
	}

	return &role, nil
}

func (r *PgRoleRepository) AssignRole(
	ctx context.Context,
	userID string,
	roleID string,
	organizationID *string,
	grantedBy *string,
	expiresAt *time.Time,
) error {

	const q = `
        INSERT INTO auth.user_roles (
            user_id,
            role_id,
            organization_id,
            granted_by,
            expires_at
        )
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id, role_id, organization_id)
        DO NOTHING;
    `

	_, err := r.db.ExecContext(
		ctx,
		q,
		userID,
		roleID,
		organizationID,
		grantedBy,
		expiresAt,
	)
	return err
}
