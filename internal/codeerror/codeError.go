package codeerror

import (
	"net/http"
)

// CodeError type
type CodeError struct {
	code      int
	qualifier string
	message   string
}

func (e CodeError) Error() string {
	return e.message
}

// Code function
func (e CodeError) Code() int {
	return e.code
}

// Qualifier function
func (e CodeError) Qualifier() string {
	return e.qualifier
}

// New function
func New(code int, qualifier string, text string) *CodeError {
	return &CodeError{code, qualifier, text}
}

// NewInternalServerError function
func NewInternalServerError(text string) *CodeError {
	return &CodeError{http.StatusInternalServerError, "", text}
}

// NewBadRequest function
func NewBadRequest(text string) *CodeError {
	return &CodeError{http.StatusBadRequest, "", text}
}

// NewNotFound function
func NewNotFound(text string) *CodeError {
	return &CodeError{http.StatusNotFound, "", text}
}

// NewUnauthorized function
func NewUnauthorized(text string) *CodeError {
	return &CodeError{http.StatusUnauthorized, "basic", text}
}

// NewUnauthorizedJWTExpired function
func NewUnauthorizedJWTExpired(text string) *CodeError {
	return &CodeError{http.StatusUnauthorized, "jwt-expired", text}
}
