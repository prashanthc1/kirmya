package common

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Error   *AppError   `json:"error"`
	Data    interface{} `json:"data"`
	Meta    Meta        `json:"meta"`
}

var (
	ErrValidationFailed = &AppError{
		Code:   "validation_failed",
		Status: http.StatusBadRequest,
	}
	ErrNotFound = &AppError{
		Code:    "not_found",
		Message: "Resource not found",
		Status:  http.StatusNotFound,
	}
	ErrUnauthorized = &AppError{
		Code:    "unauthorized",
		Message: "Invalid credentials",
		Status:  http.StatusUnauthorized,
	}
	ErrForbidden = &AppError{
		Code:    "forbidden",
		Message: "Access denied",
		Status:  http.StatusForbidden,
	}
	ErrConflict = &AppError{
		Code:    "conflict",
		Message: "Resource already exists",
		Status:  http.StatusConflict,
	}
	ErrInternalServer = &AppError{
		Code:    "internal_error",
		Message: "Internal server error",
		Status:  http.StatusInternalServerError,
	}
)

func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    "validation_failed",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    "not_found",
		Message: message,
		Status:  http.StatusNotFound,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    "unauthorized",
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    "forbidden",
		Message: message,
		Status:  http.StatusForbidden,
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    "conflict",
		Message: message,
		Status:  http.StatusConflict,
	}
}

func NewInternalError(message string) *AppError {
	return &AppError{
		Code:    "internal_error",
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func WriteError(w http.ResponseWriter, err *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)

	response := ErrorResponse{
		Success: false,
		Error:   err,
		Data:    nil,
		Meta:    getMeta(),
	}

	json.NewEncoder(w).Encode(response)
}

func WriteValidationError(w http.ResponseWriter, message string) {
	WriteError(w, NewValidationError(message))
}

func WriteNotFoundError(w http.ResponseWriter, message string) {
	WriteError(w, NewNotFoundError(message))
}

func WriteUnauthorizedError(w http.ResponseWriter, message string) {
	WriteError(w, NewUnauthorizedError(message))
}

func WriteForbiddenError(w http.ResponseWriter, message string) {
	WriteError(w, NewForbiddenError(message))
}

func WriteInternalError(w http.ResponseWriter, message string) {
	WriteError(w, NewInternalError(message))
}
