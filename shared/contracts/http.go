package contracts

import (
	"encoding/json"

	"github.com/iamonah/rideshare/shared/errs"
)

// APIResponse is the response structure for the API.
type APIResponse struct {
	Data  json.RawMessage `json:"data,omitempty"`
	Error *APIError       `json:"error,omitempty"`
}

type APIError struct {
	Code    string       `json:"code"`
	Message string       `json:"message,omitempty"`
	Fields  []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewAPIError(err *errs.AppError) *APIError {
	if err == nil {
		return nil
	}

	apiErr := &APIError{
		Code:    err.Code.String(),
		Message: err.Message,
	}

	if len(err.Fields) == 0 {
		return apiErr
	}

	apiErr.Fields = make([]FieldError, 0, len(err.Fields))
	for _, fieldErr := range err.Fields {
		apiErr.Fields = append(apiErr.Fields, FieldError{
			Field:   fieldErr.Field,
			Message: fieldErr.Message,
		})
	}

	return apiErr
}
