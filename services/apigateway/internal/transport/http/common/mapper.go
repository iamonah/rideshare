package httpcommon

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/iamonah/rideshare/shared/errs"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func WriteUpstreamGRPCError(w http.ResponseWriter, serviceName string, err error) error {
	return WriteAPIError(w, appErrorFromUpstreamGRPC(serviceName, err))
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
		apiErr.Message = unifiedUpstreamMessage()
		return apiErr
	}

	st, ok := status.FromError(err)
	if !ok {
		return apiErr
	}

	apiErr.Code = errs.FromGRPCCode(st.Code())
	if message := unifiedUpstreamMessageForCode(st.Code()); message != "" {
		apiErr.Message = message
	} else if st.Message() != "" {
		apiErr.Message = st.Message()
	}
	apiErr.Fields = apiFieldErrors(st)

	return apiErr
}

func unifiedUpstreamMessage() string {
	return "Service temporarily unavailable, please try again later"
}

func unifiedUpstreamMessageForCode(code codes.Code) string {
	switch code {
	case codes.Canceled, codes.DeadlineExceeded, codes.Unknown, codes.Internal, codes.Unavailable, codes.DataLoss:
		return unifiedUpstreamMessage()
	default:
		return ""
	}
}
