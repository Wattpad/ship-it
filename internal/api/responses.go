package api

import (
	"encoding/json"
	"net/http"
)

func Success200(w http.ResponseWriter, v interface{}) error {
	return jsonWrite(w, v, http.StatusOK)
}

func jsonWrite(w http.ResponseWriter, v interface{}, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

// Error is a structured JSON error response from the API
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewError(code int, err error) Error {
	return Error{
		Code:    code,
		Message: err.Error(),
	}
}

func ErrorJSON(w http.ResponseWriter, err Error) error {
	return jsonWrite(w, err, err.Code)
}

func Error400(w http.ResponseWriter, err error) error {
	return ErrorJSON(w, NewError(http.StatusBadRequest, err))
}

func Error404(w http.ResponseWriter, err error) error {
	return ErrorJSON(w, NewError(http.StatusNotFound, err))
}

func Error500(w http.ResponseWriter, err error) error {
	return ErrorJSON(w, NewError(http.StatusInternalServerError, err))
}
