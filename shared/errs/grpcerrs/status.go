package grpcerrs

import (
	"errors"

	sharederrs "github.com/iamonah/rideshare/shared/errs"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToStatus(err error) error {
	if err == nil {
		return nil
	}

	var appErr *sharederrs.AppError
	if errors.As(err, &appErr) && appErr != nil {
		if len(appErr.Fields) > 0 {
			return statusWithFieldDetails(appErr.Code.GRPCStatus(), appErr.Error(), fieldViolations(appErr.Fields))
		}

		return status.Error(appErr.Code.GRPCStatus(), clientMessage(appErr))
	}

	var fields sharederrs.FieldErrors
	if errors.As(err, &fields) && len(fields) > 0 {
		return statusWithFieldDetails(sharederrs.InvalidArgument.GRPCStatus(), "validation failed", fieldViolations(fields))
	}

	return status.Error(sharederrs.Internal.GRPCStatus(), "internal service error")
}

func clientMessage(appErr *sharederrs.AppError) string {
	if appErr == nil {
		return "internal service error"
	}

	switch appErr.Code {
	case sharederrs.Internal, sharederrs.Unknown, sharederrs.DataLoss, sharederrs.InternalOnlyLog:
		return "internal service error"
	default:
		return appErr.Error()
	}
}

func statusWithFieldDetails(code codes.Code, message string, violations []*errdetails.BadRequest_FieldViolation) error {
	st := status.New(code, message)
	if len(violations) == 0 {
		return st.Err()
	}

	withDetails, err := st.WithDetails(&errdetails.BadRequest{
		FieldViolations: violations,
	})
	if err != nil {
		return st.Err()
	}

	return withDetails.Err()
}

func fieldViolations(fields sharederrs.FieldErrors) []*errdetails.BadRequest_FieldViolation {
	if len(fields) == 0 {
		return nil
	}

	violations := make([]*errdetails.BadRequest_FieldViolation, 0, len(fields))
	for _, fieldErr := range fields {
		violations = append(violations, &errdetails.BadRequest_FieldViolation{
			Field:       fieldErr.Field,
			Description: fieldErr.Message,
		})
	}

	return violations
}
