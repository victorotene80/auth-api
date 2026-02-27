package response

import "time"

type CreateUserResponse struct {
	UserID                string    `json:"user_id"`
	Email                 string    `json:"email"`
	FirstName             string    `json:"first_name"`
	LastName              string    `json:"last_name"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}