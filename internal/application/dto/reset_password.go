package dto

type ResetPasswordDTO struct {
	Email      string `json:"email"`
	ResetToken string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}
