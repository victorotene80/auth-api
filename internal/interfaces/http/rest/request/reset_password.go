package request

type ResetPasswordRequest struct {
	ResetToken  string `json:"reset_token" binding:"required,len=64"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
	IPAddress   string `json:"ip_address,omitempty" binding:"omitempty,ip"`
	UserAgent   string `json:"user_agent,omitempty" binding:"omitempty"`
	DeviceID    string `json:"device_id,omitempty" binding:"omitempty,uuid4"`
}