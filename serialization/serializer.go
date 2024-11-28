
// File: serializer.go

package serialization

// Serializer defines the interface for object serialization.
type Serializer interface {
	Marshal(value interface{}) ([]byte, error)
	Unmarshal(data []byte, value interface{}) error
}
