package alfred

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Items   []Item `json:"items"`
	Keyword string `json:"-"`
}

func (r Result) ToJSON() ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}
	return data, nil
}
