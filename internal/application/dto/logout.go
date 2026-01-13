package dto

import "time"

type LogoutDTO struct {
	SessionID string
	Status    string
	Reason    string
	Time      time.Time
}
