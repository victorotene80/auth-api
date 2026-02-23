package request

type LogoutRequest struct {
	UserID    string `json:"user_id" binding:"required,uuid4"`
	SessionID string `json:"session_id" binding:"required,uuid4"`
	DeviceID  string `json:"device_id,omitempty" binding:"omitempty,uuid4"`
	IPAddress string `json:"ip_address,omitempty" binding:"omitempty,ip"`
	UserAgent string `json:"user_agent,omitempty" binding:"omitempty"`
	Reason    string `json:"reason,omitempty" binding:"omitempty"`
}
