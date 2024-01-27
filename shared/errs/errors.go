package errs

import (
	"encoding/json"
	"net/http"

	"github.com/esmailemami/chess/shared/consts"
	"github.com/esmailemami/chess/shared/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type AppError interface {
	Error() string
	Msg(string) AppError
	WithError(error) AppError
	GetStatusCode() int
	getLogError() error
	getMessage() string
}

func ErrorHandler(w http.ResponseWriter, err error) {
	switch appErr := err.(type) {
	case *Error:
		sendJSONResponse(w, appErr.GetStatusCode(), appErr)
	case *ValidationError:
		sendJSONResponse(w, appErr.GetStatusCode(), appErr)
	default:
		appErr = InternalServerErr().WithError(err)
		sendJSONResponse(w, http.StatusInternalServerError, appErr)
	}
}

// sendJSONResponse sends a JSON response with the specified status code
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
	}
}

type ValidationError struct {
	Message  string `json:"message"`
	status   int    `json:"-"`
	logError error
	Errs     map[string]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Msg(msg string) AppError {
	e.Message = msg

	return e
}

func (e *ValidationError) WithError(logError error) AppError {
	e.logError = logError
	log(e)
	return e
}

func (e *ValidationError) GetStatusCode() int {
	return e.status
}

func (e *ValidationError) getLogError() error {
	return e.logError
}

func (e *ValidationError) getMessage() string {
	return e.Message
}

func ValidationErr(err error) AppError {
	e := &ValidationError{
		Message: consts.ValidationError,
		status:  http.StatusUnprocessableEntity,
	}

	if valErr, ok := err.(validation.Errors); ok {
		e.Errs = make(map[string]string, len(valErr))
		for k, v := range valErr {
			e.Errs[k] = v.Error()
		}
	}

	return e
}

// ------------------------------------------------------

type Error struct {
	Message  string `json:"message"`
	status   int    `json:"-"`
	logError error
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Msg(msg string) AppError {
	e.Message = msg

	return e
}

func (e *Error) WithError(logError error) AppError {
	e.logError = logError
	log(e)
	return e
}

func (e *Error) GetStatusCode() int {
	return e.status
}

func (e *Error) getLogError() error {
	return e.logError
}

func (e *Error) getMessage() string {
	return e.Message
}

func InternalServerErr() AppError {
	e := &Error{
		Message: consts.InternalServerError,
		status:  http.StatusInternalServerError,
	}

	return e
}

func UnAuthorizedErr() AppError {
	e := &Error{
		Message: consts.UnauthorizedError,
		status:  http.StatusUnauthorized,
	}

	return e
}

func NotFoundErr() AppError {
	e := &Error{
		Message: consts.RecordNotFound,
		status:  http.StatusNotFound,
	}

	return e
}

func BadRequestErr() AppError {
	e := &Error{
		Message: consts.BadRequest,
		status:  http.StatusBadRequest,
	}

	return e
}

func log(e AppError) {
	logging.ErrorE(e.getMessage(), e.getLogError(), "status", e.GetStatusCode())
}
