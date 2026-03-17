package errs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(jsonTagName)
	return v
}

func Validate(data any) error {
	return ValidateStruct(data)
}

func ValidateStruct(data any) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fieldErrors := NewFieldErrors()
	for _, verr := range validationErrors {
		fieldErrors.AddMessage(fieldPath(verr), validateMessage(verr))
	}

	return fieldErrors
}

func jsonTagName(field reflect.StructField) string {
	name := strings.Split(field.Tag.Get("json"), ",")[0]
	if name == "-" {
		return ""
	}
	if name != "" {
		return name
	}
	return field.Name
}

func fieldPath(verr validator.FieldError) string {
	namespace := verr.Namespace()
	if namespace == "" {
		return verr.Field()
	}

	if dot := strings.Index(namespace, "."); dot >= 0 {
		return namespace[dot+1:]
	}

	return namespace
}

func validateMessage(verr validator.FieldError) string {
	switch verr.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s", verr.Param())
	case "max":
		return fmt.Sprintf("must be at most %s", verr.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s", verr.Param())
	case "oneof":
		return fmt.Sprintf("must be one of %s", verr.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", verr.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", verr.Param())
	case "lt":
		return fmt.Sprintf("must be less than %s", verr.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", verr.Param())
	default:
		return verr.Error()
	}
}
