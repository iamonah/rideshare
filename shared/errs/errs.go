package errs

import (
	"errors"
	"fmt"
	"runtime"
)

type AppError struct {
	Code     ErrCode     
	Message  string      
	Err      error       
	Fields   FieldErrors 
	FuncName string      
	FileName string      
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

func New(code ErrCode, err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) && appErr != nil && len(appErr.Fields) > 0 {
		return newAppError(1, code, appErr.Message, nil, appErr.Fields)
	}

	var fieldErrs FieldErrors
	if errors.As(err, &fieldErrs) && len(fieldErrs) > 0 {
		return newAppError(1, code, "validation failed", nil, fieldErrs)
	}

	return newAppError(1, code, err.Error(), err, nil)
}

func Newf(code ErrCode, err error, format string, v ...any) *AppError {
	if err == nil {
		err = fmt.Errorf(format, v...)
		return newAppError(1, code, err.Error(), err, nil)
	}
	return newAppError(1, code, fmt.Sprintf(format, v...), err, nil)
}

func Validation(fields FieldErrors) *AppError {
	return validationMessage(2, "validation failed", fields)
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
