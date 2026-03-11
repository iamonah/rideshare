package errs

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/iamonah/verrors"
)

type AppErr struct {
	Kind     string              `json:"error"`
	Code     ErrCode             `json:"status_code"`
	Message  string              `json:"message,omitempty"`
	Fields   verrors.FieldErrors `json:"fields,omitempty"`
	Cause    error               `json:"-"`
	FuncName string              `json:"-"`
	FileName string              `json:"-"`
}

func (e *AppErr) Error() string {
	if e == nil {
		return ""
	}
	if e.Message != "" {
		return e.Message
	}
	return e.Code.String()
}

func (e *AppErr) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func New(code ErrCode, err error) *AppErr {
	return newAppErr(1, code, err)
}

func Newf(code ErrCode, format string, v ...any) *AppErr {
	return newAppErr(1, code, fmt.Errorf(format, v...))
}

func newAppErr(skip int, code ErrCode, err error) *AppErr {
	code = resolveCode(code, err)
	pc, filename, line, ok := runtime.Caller(skip + 1)
	appErr := &AppErr{
		Kind: code.String(),
		Code: code,
	}
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			appErr.FuncName = fn.Name()
		}
		appErr.FileName = fmt.Sprintf("%s:%d", filename, line)
	}

	if err == nil {
		appErr.Message = code.String()
		return appErr
	}

	var wrapped *AppErr
	if errors.As(err, &wrapped) && wrapped != nil {
		appErr.Message = wrapped.Message
		appErr.Fields = wrapped.Fields
		appErr.Cause = wrapped.Cause
		if appErr.Message == "" {
			appErr.Message = wrapped.Error()
		}
		if appErr.Cause == nil {
			appErr.Cause = err
		}
		return appErr
	}

	var fields verrors.FieldErrors
	if errors.As(err, &fields) && len(fields) > 0 {
		appErr.Message = "validation failed"
		appErr.Fields = fields
		appErr.Cause = err
		return appErr
	}

	var domainErr *DomainError
	if errors.As(err, &domainErr) && domainErr != nil {
		appErr.Message = domainErr.Error()
		appErr.Cause = domainErr.Err
		if appErr.Cause == nil {
			appErr.Cause = err
		}
		return appErr
	}

	appErr.Message = err.Error()
	appErr.Cause = err
	return appErr
}

func resolveCode(code ErrCode, err error) ErrCode {
	if err == nil || (code != None && code != Unknown) {
		return code
	}

	var wrapped *AppErr
	if errors.As(err, &wrapped) && wrapped != nil && wrapped.Code != None {
		return wrapped.Code
	}

	var domainErr *DomainError
	if errors.As(err, &domainErr) && domainErr != nil && domainErr.Code != None {
		return domainErr.Code
	}

	var fields verrors.FieldErrors
	if errors.As(err, &fields) && len(fields) > 0 {
		return InvalidArgument
	}

	if code == None {
		return Unknown
	}

	return code
}

// errs we expect are bound to happen
type DomainError struct {
	Err  error
	Code ErrCode
}

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code.String()
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func NewDomainError(code ErrCode, err error) *DomainError {
	return &DomainError{Code: code, Err: err}
}

func IsDomainError(err error) (*DomainError, bool) {
	var dError *DomainError
	if errors.As(err, &dError) {
		return dError, true
	}
	return nil, false
}
