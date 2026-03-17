package apigateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iamonah/rideshare/shared/errs"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

func writeUpstreamGRPCError(w http.ResponseWriter, serviceName string, err error) {
	writeAPIError(w, appErrorFromUpstreamGRPC(serviceName, err))
}

func apiFieldErrors(st *status.Status) errs.FieldErrors {
	if st == nil {
		return nil
	}

	fields := make(errs.FieldErrors, 0)
	for _, detail := range st.Details() {
		badRequest, ok := detail.(*errdetails.BadRequest)
		if !ok {
			continue
		}

		for _, violation := range badRequest.FieldViolations {
			fields = append(fields, errs.FieldError{
				Field:   violation.Field,
				Message: violation.Description,
			})
		}
	}

	if len(fields) == 0 {
		return nil
	}

	return fields
}

func appErrorFromUpstreamGRPC(serviceName string, err error) *errs.AppError {
	apiErr := errs.New(errs.Unavailable, fmt.Errorf("failed to call %s", serviceName))

	if errors.Is(err, context.DeadlineExceeded) {
		apiErr.Code = errs.DeadlineExceeded
		apiErr.Message = fmt.Sprintf("%s request timed out", serviceName)
		return apiErr
	}

	st, ok := status.FromError(err)
	if !ok {
		return apiErr
	}

	apiErr.Code = errs.FromGRPCCode(st.Code())
	if st.Message() != "" {
		apiErr.Message = st.Message()
	}
	apiErr.Fields = apiFieldErrors(st)

	return apiErr
}
