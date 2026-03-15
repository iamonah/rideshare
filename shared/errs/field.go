package errs

import "encoding/json"

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrors []FieldError

func NewFieldErrors() FieldErrors {
	return FieldErrors{}
}

func (fe *FieldErrors) Add(field string, err error) {
	if err == nil {
		return
	}
	fe.AddMessage(field, err.Error())
}

func (fe *FieldErrors) AddMessage(field, message string) {
	*fe = append(*fe, FieldError{
		Field:   field,
		Message: message,
	})
}

func (fe FieldErrors) ToError() error {
	if len(fe) == 0 {
		return nil
	}
	return validationMessage(2, "validation failed", fe)
}

func (fe FieldErrors) ToErrorWithMessage(message string) error {
	if len(fe) == 0 {
		return nil
	}
	return validationMessage(2, message, fe)
}

func (fe FieldErrors) Error() string {
	d, err := json.Marshal(fe)
	if err != nil {
		return err.Error()
	}
	return string(d)
}
