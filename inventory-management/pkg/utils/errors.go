package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrCodeInvalidInput = "INVALID_INPUT"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeInternal     = "INTERNAL_ERROR"
)

type AppError struct {
	Code    string
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func NewInvalidInputError(message string, cause error) *AppError {
	return &AppError{Code: ErrCodeInvalidInput, Message: message, Cause: cause}
}

func NewNotFoundError(message string, cause error) *AppError {
	return &AppError{Code: ErrCodeNotFound, Message: message, Cause: cause}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{Code: ErrCodeInternal, Message: message, Cause: cause}
}

func HTTPStatus(err error) int {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		return http.StatusInternalServerError
	}

	switch appErr.Code {
	case ErrCodeInvalidInput:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func PublicMessage(err error) string {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		return "terjadi kesalahan internal"
	}
	return appErr.Message
}

// HandleError dipakai semua handler agar mapping error ke HTTP konsisten.
func HandleError(c *gin.Context, err error) {
	Error(c, HTTPStatus(err), PublicMessage(err), nil)
}
