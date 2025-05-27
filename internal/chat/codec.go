package chat

import (
	"encoding/json"
	"fmt"
)

func Decode[T any](data []byte) (*T, error) {
	var out T
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("Error unmarshaling data into type %T: %w", out, err)
	}
	return &out, nil
}

func Encode(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling %T: %w", v, err)
	}
	return data, nil
}
