package request

type CreateUserRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=8"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name,omitempty" binding:"omitempty"`
	DeviceID   string `validate:"required,uuid4"`
	IPAddress  string `validate:"required,ip"`
	UserAgent  string `validate:"required"`
}
