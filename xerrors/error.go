package xerrors

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
)

type Error struct {
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	switch {
	case e == nil:
		return ""
	case e.Code == "" && e.Message == "" && e.Err != nil:
		return e.Err.Error()
	case e.Code == "":
		if e.Err == nil {
			return e.Message
		}
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	case e.Message == "":
		if e.Err == nil {
			return e.Code
		}
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	default:
		if e.Err == nil {
			return fmt.Sprintf("%s: %s", e.Code, e.Message)
		}
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func New(code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(err error, code, message string) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

func CodeOf(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}

func IsCode(err error, code string) bool {
	return CodeOf(err) == code
}

var (
	statusMu     sync.RWMutex
	statusByCode = map[string]int{
		"INVALID_ARGUMENT":  http.StatusBadRequest,
		"UNAUTHORIZED":      http.StatusUnauthorized,
		"FORBIDDEN":         http.StatusForbidden,
		"NOT_FOUND":         http.StatusNotFound,
		"CONFLICT":          http.StatusConflict,
		"TOO_MANY_REQUESTS": http.StatusTooManyRequests,
	}
)

func RegisterHTTPStatus(code string, status int) {
	if code == "" || status <= 0 {
		return
	}
	statusMu.Lock()
	statusByCode[code] = status
	statusMu.Unlock()
}

func HTTPStatus(err error) int {
	code := CodeOf(err)
	if code == "" {
		return http.StatusInternalServerError
	}
	statusMu.RLock()
	status, ok := statusByCode[code]
	statusMu.RUnlock()
	if !ok {
		return http.StatusInternalServerError
	}
	return status
}
