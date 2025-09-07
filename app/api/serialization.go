package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

var DeserializationError = errors.New("invalid JSON format")

func ParseRequestBodyJSON[T any](r *http.Request) (T, error) {
	var payload T

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("Error closing request body in parseRequestBodyJSON, error: %v\n", err)
		}
	}(r.Body)

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Invalid JSON format for parsing, error: %v\n", err)
		return payload, DeserializationError
	}

	return payload, nil
}
