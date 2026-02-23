package request

type LockAccountRequest struct {
	Reason string `json:"reason,omitempty"`
}