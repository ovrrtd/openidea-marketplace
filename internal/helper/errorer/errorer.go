package errorer

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrBadRequest       = errors.New("bad request")
	ErrNotFound         = errors.New("not found")
	ErrInternalServer   = errors.New("internal server error")
	ErrInternalDatabase = errors.New("internal database error")
	ErrEmailExist       = errors.New("email already exist")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
)

func ErrInputRequest(err error) error {
	return fmt.Errorf("input request error: %s", err.Error())
}

func HTTPCodeFromError(err error) int {
	if err == ErrBadRequest {
		return http.StatusBadRequest
	} else if err == ErrNotFound {
		return http.StatusNotFound
	} else if err == ErrInternalServer {
		return http.StatusInternalServerError
	} else if err == ErrEmailExist {
		return http.StatusBadRequest
	} else if strings.HasPrefix(err.Error(), "input request error:") {
		return http.StatusBadRequest
	} else if err == ErrUnauthorized {
		return http.StatusUnauthorized
	} else {
		return http.StatusInternalServerError
	}
}
