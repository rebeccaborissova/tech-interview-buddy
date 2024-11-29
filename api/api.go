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

	Session string

	Message string
}

// Sign-Up parameters
type SignUpParams struct {
	Username string
	Authorization string
	FirstName string
	LastName string
	DSA bool
	Year int
	Description string
}

// Sign-Up reponse
type SignUpResponse struct {
	// Success code
	Code int

	Username string

	Message string
}

// Generic response
type SimpleResponse struct {
	// Success code
	Code int

	Message string
}

// User info response
type UserInfoResponse struct {
	// Success code
	Code int

	// User attributes
	Username string
	FirstName string
	LastName  string
	InvitedBy string
	TakenDSA bool
	Year   int
	Description string
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
