package common

import "encoding/json"

type ErrorMessage struct {
	Message string `json:"message"`
}

func (errorMessage *ErrorMessage) ToBytes() []byte {
	bytes, _ := json.Marshal(errorMessage)
	return bytes
}

func NewErrorMessage(err error) *ErrorMessage {
	return &ErrorMessage{Message: err.Error()}
}

func NewErrorToJson(err error) []byte {
	return NewErrorMessage(err).ToBytes()
}
