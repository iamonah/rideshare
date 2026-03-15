package errs

import (
	"fmt"
	"runtime"
)

type AppError struct {
	Code     ErrCode     `json:"code"`
	Message  string      `json:"message,omitempty"`
	Err      error       `json:"-"`
	Fields   FieldErrors `json:"fields,omitempty"`
	FuncName string      `json:"-"`
	FileName string      `json:"-"`
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	if len(e.Fields) > 0 {
		return "validation failed"
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code.String()
}

func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func New(code ErrCode, message string) *AppError {
	return newAppError(1, code, message, nil, nil)
}

func Newf(code ErrCode, format string, v ...any) *AppError {
	return newAppError(1, code, fmt.Sprintf(format, v...), nil, nil)
}

func Wrap(code ErrCode, message string, err error) *AppError {
	return newAppError(1, code, message, err, nil)
}

func Wrapf(code ErrCode, err error, format string, v ...any) *AppError {
	return newAppError(1, code, fmt.Sprintf(format, v...), err, nil)
}

func Validation(fields FieldErrors) *AppError {
	return validationMessage(2, "validation failed", fields)
}

func ValidationMessage(message string, fields FieldErrors) *AppError {
	return validationMessage(2, message, fields)
}

func validationMessage(skip int, message string, fields FieldErrors) *AppError {
	if len(fields) == 0 {
		return nil
	}

	copied := append(FieldErrors(nil), fields...)
	if message == "" {
		message = "validation failed"
	}

	return newAppError(skip, InvalidArgument, message, nil, copied)
}

func newAppError(skip int, code ErrCode, message string, err error, fields FieldErrors) *AppError {
	appErr := &AppError{
		Code:    code,
		Message: message,
		Err:     err,
		Fields:  fields,
	}

	pc, filename, line, ok := runtime.Caller(skip + 1)
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			appErr.FuncName = fn.Name()
		}
		appErr.FileName = fmt.Sprintf("%s:%d", filename, line)
	}

	return appErr
}

func NewValidation(fields FieldErrors) *AppError {
	return validationMessage(2, "validation failed", fields)
}

func NewValidationMessage(message string, fields FieldErrors) *AppError {
	return validationMessage(2, message, fields)
}
