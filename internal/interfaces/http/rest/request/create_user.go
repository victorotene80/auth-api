package request

type CreateUserRequest struct {
    Email       string  `json:"email"`
    Password    string  `json:"password"`
    FirstName   string  `json:"firstName"`
    LastName    string  `json:"lastName"`
    MiddleName  *string `json:"middleName,omitempty"`
    Phone       *string `json:"phone,omitempty"`
    Locale      *string `json:"locale,omitempty"`
    Timezone    *string `json:"timezone,omitempty"`
    AcceptTerms bool    `json:"acceptTerms"`
}