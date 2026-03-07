package messaging

import "encoding/json"

func MustJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func JSON(v any) ([]byte, error) {
	return json.Marshal(v)
}