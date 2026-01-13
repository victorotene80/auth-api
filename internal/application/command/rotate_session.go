package command

type RotateSessionCommand struct {
	RawToken   string // Presented by client (raw token)
	IPAddress  string // Optional: for logging / analytics
	UserAgent  string // Optional: for logging / analytics
	DeviceID   string
	//RotationID string // Unique idempotency token for rotation
	// NewExpiresAt *time.Time // Optional, only if you want to extend session
}
