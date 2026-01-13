package dto

import "time"

type ChangePasswordDTO struct {
	UserID    string
	Success   bool
	ChangedAt time.Time
}
