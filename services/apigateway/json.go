package apigateway

import (
	"encoding/json"
	"net/http"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/errs"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func writeAPIError(w http.ResponseWriter, err *errs.AppError) error {
	return writeJSON(w, err.Code.HTTPStatus(), contracts.APIResponse{
		Error: contracts.NewAPIError(err),
	})
}
