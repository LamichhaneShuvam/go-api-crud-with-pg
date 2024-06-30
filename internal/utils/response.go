package utils

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Success struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func writeSuccess(w http.ResponseWriter, data any, message string, code int) {
	response := Success{
		Code:    code,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

func writeError(w http.ResponseWriter, message string, code int) {
	response := Error{
		Code:    code,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

var (
	//* Error Vars
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "Something went wrong", http.StatusInternalServerError)
	}
	NotFoundErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusNotFound)
	}

	//* Success vars
	OkResponseHandler = func(w http.ResponseWriter, data any, message string) {
		writeSuccess(w, data, message, http.StatusOK)
	}
	CreateResponseHandler = func(w http.ResponseWriter, data any, message string) {
		writeSuccess(w, data, message, http.StatusCreated)
	}
)
