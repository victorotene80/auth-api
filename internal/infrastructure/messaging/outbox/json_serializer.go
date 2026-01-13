package outbox

import "encoding/json"

type JSONSerializer struct{}

func (s JSONSerializer) Serialize(v any) ([]byte, error) {
	return json.Marshal(v)
}
