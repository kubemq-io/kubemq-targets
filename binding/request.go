package binding

import "encoding/json"

type Request struct {
	Binding string          `json:"binding"`
	Payload json.RawMessage `json:"payload"`
}
