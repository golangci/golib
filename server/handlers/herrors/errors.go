package herrors

import (
	"fmt"
	"net/http"
)

type HTTPError interface {
	error
	Code() int
}

type httpError struct {
	code    int
	message string
}

func (e httpError) Code() int {
	return e.code
}

func (e httpError) Error() string {
	return e.message
}

func New400Errorf(format string, args ...interface{}) HTTPError {
	return httpError{
		code:    http.StatusBadRequest,
		message: fmt.Sprintf(format, args...),
	}
}

func New404Errorf(format string, args ...interface{}) HTTPError {
	return httpError{
		code:    http.StatusNotFound,
		message: fmt.Sprintf(format, args...),
	}
}

func New403Errorf(format string, args ...interface{}) HTTPError {
	return httpError{
		code:    http.StatusForbidden,
		message: fmt.Sprintf(format, args...),
	}
}

func New(err error, format string, args ...interface{}) HTTPError {
	newMessage := fmt.Sprintf(format, args...) + ": " + err.Error()

	herr, ok := err.(httpError)
	if ok {
		return httpError{
			code:    herr.code,
			message: newMessage,
		}
	}

	return httpError{
		code:    http.StatusInternalServerError,
		message: newMessage,
	}
}
