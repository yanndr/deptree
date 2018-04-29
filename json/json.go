package json

import (
	"encoding/json"
	"fmt"
	"os"
)

// DecodeFromFile reads the next JSON-encoded value from its
// input and stores it in the value pointed to by v.
func DecodeFromFile(v interface{}, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening the file %s , %s", path, err)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(v)
	if err != nil {
		return fmt.Errorf("error decoding the the file %s, %s", file.Name(), err)
	}

	return nil
}
