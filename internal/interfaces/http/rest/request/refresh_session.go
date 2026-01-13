package request

type RefreshSessionRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,len=64"`
	IPAddress    string `json:"ip_address,omitempty" binding:"omitempty,ip"`
	UserAgent    string `json:"user_agent,omitempty" binding:"omitempty"`
	DeviceID     string `json:"device_id,omitempty" binding:"omitempty,uuid4"`
}