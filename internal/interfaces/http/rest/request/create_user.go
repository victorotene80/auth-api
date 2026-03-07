package request

type CreateUserRequest struct {
	Email       string  `json:"email"       validate:"required,email"`
	Password    string  `json:"password"    validate:"required,min=8"`
	FirstName   string  `json:"firstName"   validate:"required,min=1,max=100"`
	LastName    string  `json:"lastName"    validate:"required,min=1,max=100"`
	MiddleName  *string `json:"middleName,omitempty" validate:"omitempty,min=1,max=100"`
	Phone       *string `json:"phone,omitempty"     validate:"omitempty,e164"`
	Locale      *string `json:"locale,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	AcceptTerms bool    `json:"acceptTerms" validate:"required,eq=true"`
}
