package integration

import (
	"encoding/json"
	"io"
)

func ReaderToStruct(body io.ReadCloser, value any) error {
	bytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, value)
	if err != nil {
		return err
	}
	return nil
}
