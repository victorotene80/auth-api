package response

type VerifyMFAResponse struct {
	UserID string `json:"user_id"`
	Method string `json:"method"`
	Status string `json:"status"` 
}