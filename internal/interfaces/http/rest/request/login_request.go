package request

type LoginRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	IPAddress string `json:"ip_address,omitempty" binding:"omitempty,ip"`
	UserAgent string `json:"user_agent,omitempty" binding:"omitempty"`
	DeviceID  string `json:"device_id,omitempty" binding:"omitempty,uuid4"`
}
