package constant

import (
	"errors"
	"net/http"
)

const (
	ErrSQLUniqueViolation = "23505"
)

var (
	ErrEmailAlreadyRegistered = ErrWithCode{HTTPStatusCode: http.StatusConflict, Message: "email already registered"}
	ErrEmailOrPasswordInvalid = ErrWithCode{HTTPStatusCode: http.StatusBadRequest, Message: "email or password invalid"}
	ErrInvalidUUID            = errors.New("invalid uuid length or format")
	ErrUnauthorizedAccess     = ErrWithCode{HTTPStatusCode: http.StatusUnauthorized, Message: "unauthorized access"}
)

type ErrWithCode struct {
	HTTPStatusCode int    `json:"-"`
	Message        string `json:"message"`
	Field          string `json:"field,omitempty"`
}

func (e ErrWithCode) Error() string {
	return e.Message
}

type ErrsWithCode []ErrWithCode

func (e ErrsWithCode) Error() string {
	message := ""
	for _, err := range e {
		message += err.Message + ", "
	}

	// remove last comma
	message = message[:len(message)-2]

	return message
}
