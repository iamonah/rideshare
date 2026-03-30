package httpcommon

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/iamonah/rideshare/shared/contracts"
	"github.com/iamonah/rideshare/shared/errs"
)

func ReadJSON[T any](r *http.Request, dest *T) error {
	if dest == nil {
		panic("dest must be a pointer")
	}

	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		if fallbackErr := writeFallbackInternalError(w); fallbackErr != nil {
			return errors.Join(fmt.Errorf("failed to marshal JSON response: %w", err), fallbackErr)
		}
		return fmt.Errorf("failed to marshal JSON response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(append(payload, '\n'))

	return nil
}

func writeFallbackInternalError(w http.ResponseWriter) error {
	fallbackPayload := []byte(`{"error":{"code":"internal","message":"internal service error"}}` + "\n")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	w.Write(fallbackPayload)
	return nil
}

func WriteAPIError(w http.ResponseWriter, err *errs.AppError) error {
	if err == nil {
		return WriteJSON(w, http.StatusInternalServerError, contracts.APIResponse{
			Error: contracts.NewAPIError(errs.New(errs.Internal, errors.New("internal service error"))),
		})
	}

	return WriteJSON(w, err.Code.HTTPStatus(), contracts.APIResponse{
		Error: contracts.NewAPIError(err),
	})
}
