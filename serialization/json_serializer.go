
// File: json_serializer.go

package serialization

import (
	"encoding/json"
)

// JSONSerializer serializes objects into JSON format.
type JSONSerializer struct{}

// Marshal converts an object into JSON.
func (s *JSONSerializer) Marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// Unmarshal converts JSON data into an object.
func (s *JSONSerializer) Unmarshal(data []byte, value interface{}) error {
	return json.Unmarshal(data, value)
}
