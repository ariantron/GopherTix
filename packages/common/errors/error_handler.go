package errors

import (
	"errors"
	"net/http"
)

type AppError struct {
	err        error
	statusCode int
	code       string
	details    map[string]interface{}
}

type ErrorOption func(*AppError)

func WithStatus(status int) ErrorOption {
	return func(e *AppError) {
		e.statusCode = status
	}
}

func WithCode(code string) ErrorOption {
	return func(e *AppError) {
		e.code = code
	}
}

func WithDetails(details map[string]interface{}) ErrorOption {
	return func(e *AppError) {
		if e.details == nil {
			e.details = make(map[string]interface{})
		}
		for k, v := range details {
			e.details[k] = v
		}
	}
}

func NewAppError(message string, opts ...ErrorOption) *AppError {
	appErr := &AppError{
		err:        errors.New(message),
		statusCode: http.StatusInternalServerError,
		code:       "INTERNAL_ERROR",
	}

	for _, opt := range opts {
		opt(appErr)
	}

	return appErr
}

func WrapError(err error, opts ...ErrorOption) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		for _, opt := range opts {
			opt(appErr)
		}
		return appErr
	}

	appErr = &AppError{
		err:        err,
		statusCode: http.StatusInternalServerError,
		code:       "INTERNAL_ERROR",
	}

	for _, opt := range opts {
		opt(appErr)
	}

	return appErr
}

func (e *AppError) Error() string {
	return e.err.Error()
}

func (e *AppError) Status() int {
	return e.statusCode
}

func (e *AppError) Code() string {
	return e.code
}

func (e *AppError) Details() map[string]interface{} {
	return e.details
}

func (e *AppError) Is(target error) bool {
	if target == nil {
		return false
	}

	var t *AppError
	if errors.As(target, &t) {
		return errors.Is(e.err, t.err)
	}

	return errors.Is(e.err, target)
}

func (e *AppError) Unwrap() error {
	return e.err
}

type ErrorResponse struct {
	Status  int                    `json:"status"`
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

type ErrorHandler interface {
	HandleError(err error) ErrorResponse
}

type errorHandler struct {
	defaultStatus  int
	defaultCode    string
	defaultMessage string
}

type ErrorHandlerOption func(*errorHandler)

func WithDefaultStatus(status int) ErrorHandlerOption {
	return func(h *errorHandler) {
		h.defaultStatus = status
	}
}

func WithDefaultCode(code string) ErrorHandlerOption {
	return func(h *errorHandler) {
		h.defaultCode = code
	}
}

func WithDefaultMessage(message string) ErrorHandlerOption {
	return func(h *errorHandler) {
		h.defaultMessage = message
	}
}

func NewErrorHandler(opts ...ErrorHandlerOption) ErrorHandler {
	h := &errorHandler{
		defaultStatus:  http.StatusInternalServerError,
		defaultCode:    "INTERNAL_ERROR",
		defaultMessage: "An internal error occurred",
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *errorHandler) HandleError(err error) ErrorResponse {
	if err == nil {
		return ErrorResponse{
			Status:  h.defaultStatus,
			Code:    h.defaultCode,
			Message: h.defaultMessage,
		}
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return ErrorResponse{
			Status:  appErr.Status(),
			Code:    appErr.Code(),
			Message: appErr.Error(),
			Details: appErr.Details(),
		}
	}

	return ErrorResponse{
		Status:  h.defaultStatus,
		Code:    h.defaultCode,
		Message: err.Error(),
	}
}
