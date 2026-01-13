package dto

type UpdateUserDTO struct {
	UserID    string  `json:"user_id"`
	Email     *string `json:"email,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
}
