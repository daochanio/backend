package common

import (
	"encoding/json"
	"fmt"
	"io"
)

func Unmarshal[T any](data []byte) (*T, error) {
	out := new(T)
	if err := json.Unmarshal(data, out); err != nil {
		return nil, fmt.Errorf("could not unmarshal: %w", err)
	}

	return out, nil
}

func Decode[T any](reader io.Reader) (*T, error) {
	out := new(T)
	if err := json.NewDecoder(reader).Decode(out); err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	return out, nil
}
