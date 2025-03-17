package error

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	HTTPCode     int
	ErrorDetail  string
	ErrorMessage string
	ErrorCode    string
}

func (e *CustomError) Error() string {
	return e.ErrorMessage
}

func AnyError(code int, err error) CustomError {
	return CustomError{
		HTTPCode:     code,
		ErrorMessage: err.Error(),
		ErrorCode:    fmt.Sprintf("251%v", code),
	}
}

func EmptyParam() CustomError {
	return CustomError{
		HTTPCode:     http.StatusBadRequest,
		ErrorMessage: "Empty Parameters",
		ErrorCode:    "251400",
	}
}

func InvalidParam() CustomError {
	return CustomError{
		HTTPCode:     http.StatusBadRequest,
		ErrorMessage: "Invalid Parameters",
		ErrorCode:    "251400",
	}
}

// Create ERROR

func CreateError(err error) CustomError {
	return CustomError{
		HTTPCode:     http.StatusInternalServerError,
		ErrorMessage: "Create Error",
		ErrorDetail:  err.Error(),
		ErrorCode:    "251500",
	}
}
