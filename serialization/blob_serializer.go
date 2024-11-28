
// File: blob_serializer.go

package serialization

import (
	"bytes"
	"encoding/gob"
)

// BlobSerializer serializes objects into binary blobs using gob.
type BlobSerializer struct{}

// Marshal converts an object into a binary blob.
func (s *BlobSerializer) Marshal(value interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(value)
	return buf.Bytes(), err
}

// Unmarshal converts a binary blob into an object.
func (s *BlobSerializer) Unmarshal(data []byte, value interface{}) error {
	buf := bytes.NewReader(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(value)
}
