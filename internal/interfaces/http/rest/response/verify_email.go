package response

type VerifyEmailResponse struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}
