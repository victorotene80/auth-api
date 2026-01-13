package outbox

import "fmt"

type Status string

const (
	Pending    Status = "PENDING"
	InProgress Status = "IN_PROGRESS"
	Sent       Status = "SENT"
	Failed     Status = "FAILED"
)

func (s Status) IsValid() bool {
	switch s {
	case Pending, InProgress, Sent, Failed:
		return true
	default:
		return false
	}
}

func NewStatus(value string) (Status, error) {
	s := Status(value)
	if !s.IsValid() {
		return "", fmt.Errorf("invalid outbox status: %s", value)
	}
	return s, nil

}
