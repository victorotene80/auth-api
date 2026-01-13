package request

type RevokeSessionRequest struct {
	IPAddress string `json:"ip_address,omitempty" binding:"omitempty,ip"`
	UserAgent string `json:"user_agent,omitempty" binding:"omitempty"`
	DeviceID  string `json:"device_id,omitempty" binding:"omitempty,uuid4"`
}