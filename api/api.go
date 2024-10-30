package api

import (
	"encoding/json"
	"net/http"
)

// Login parameters
type LoginParams struct {
	Username string

	Authorization string
}

// Login response
type LoginResponse struct {
	// Success code
	Code int

	Username string

	Message string
}

// Error response
type Error struct {
	// Error code
	Code int

	// Error message
	Message string
}

func writeError(writer http.ResponseWriter, message string, code int) {
	response := Error{
		Code:    code,
		Message: message,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	writer.Header().Set("Access-Control-Max-Age", "15")
	writer.WriteHeader(code)

	json.NewEncoder(writer).Encode(response)

}

var (
	RequestErrorHandler = func(writer http.ResponseWriter, err error) {
		writeError(writer, err.Error(), http.StatusBadRequest)
	}

	InternalErrorHandler = func(writer http.ResponseWriter) {
		writeError(writer, "An internal error occured.", http.StatusInternalServerError)
	}
)
