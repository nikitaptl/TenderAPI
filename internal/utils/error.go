package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Reason     string `json:"reason"`
	StatusCode int    `json:"-"`
}

func (e *ErrorResponse) Error() string {
	return e.Reason
}

func NewErrorResponse(message string, statusCode int) *ErrorResponse {
	return &ErrorResponse{message, statusCode}
}

func HandleHTTPError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var customError *ErrorResponse
	if errors.As(err, &customError) {
		w.WriteHeader(customError.StatusCode)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(err)
	log.Printf("[error]: %v", err)
}
