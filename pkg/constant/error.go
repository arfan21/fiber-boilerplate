package constant

import (
	"errors"
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrEmailAlreadyRegistered = &ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "email already registered"}
	ErrEmailOrPasswordInvalid = &ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "email or password invalid"}
	ErrInvalidUUID            = errors.New("invalid uuid length or format")
	ErrUnauthorizedAccess     = &ErrWithCode{HTTPStatusCode: http.StatusUnauthorized, Message: "unauthorized access"}
)

type ErrWithCode struct {
	HTTPStatusCode int
	Message        string
}

func (e *ErrWithCode) Error() string {
	return e.Message
}

type ErrValidation struct {
	Message string
}

func (e *ErrValidation) Error() string {
	return e.Message
}
