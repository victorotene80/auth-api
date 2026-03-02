package contracts

import (
	"context"
	"github.com/victorotene80/authentication_api/internal/application/dto"
)

type AuditLogger interface {
	Log(ctx context.Context, rec dto.AuditRecord) error
}