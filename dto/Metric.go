package dto

import "encoding/json"

type Metrics struct {
	ID    string          `json:"id"`
	MType string          `json:"type"`
	Value json.RawMessage `json:"value,omitempty"`
}
