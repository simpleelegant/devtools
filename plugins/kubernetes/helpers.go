package kubernetes

import (
	"bytes"
	"encoding/json"
)

// JSONContentType ...
const JSONContentType = "application/json"

// ToJSONReader convert data to json reader
func ToJSONReader(data interface{}) (*bytes.Reader, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
