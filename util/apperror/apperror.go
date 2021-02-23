package apperror

import (
	"fmt"
	"net/http"

	"github.com/dmitriivoitovich/heameelega/util/i18n"
)

const errorCodeValidation = http.StatusAlreadyReported + 1

type Error struct {
	HTTPCode int
	I18nKey  i18n.Key
	Err      error
	Msg      string
}

func Internal(err error, msg string) *Error {
	return &Error{
		HTTPCode: http.StatusInternalServerError,
		Err:      err,
		Msg:      msg,
	}
}

func BadRequest(err error, msg string) *Error {
	return &Error{
		HTTPCode: http.StatusBadRequest,
		Err:      err,
		Msg:      msg,
	}
}

func NotFound(msg string) *Error {
	return &Error{
		HTTPCode: http.StatusNotFound,
		Msg:      msg,
	}
}

func Validation(key i18n.Key) *Error {
	return &Error{
		HTTPCode: errorCodeValidation,
		I18nKey:  key,
	}
}

func (e *Error) Message() string {
	if e.Err == nil {
		return e.Msg
	}

	return fmt.Sprintf("%v: %s", e.Err, e.Msg)
}

func (e *Error) IsValidation() bool {
	return e.HTTPCode == errorCodeValidation
}
