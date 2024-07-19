package common

import (
	"encoding/json"
	"log"
	"net/http"
)

func SendResponse(writer http.ResponseWriter, result any) {
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(result)

	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}
