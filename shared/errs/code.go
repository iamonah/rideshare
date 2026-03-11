package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type ErrCode int

func (e ErrCode) String() string {
	if name, ok := codeNames[e]; ok {
		return name
	}
	return "unknown_error"
}

func (e ErrCode) HTTPStatus() int {
	if status, ok := httpStatusByCode[e]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func (e ErrCode) GRPCStatus() codes.Code {
	if status, ok := grpcStatusByCode[e]; ok {
		return status
	}
	return codes.Internal
}

func (e ErrCode) Error() string {
	return e.String()
}

const (
	// Keep these values stable once they are used outside this package.
	None               ErrCode = 0
	NoContent          ErrCode = 1
	Canceled           ErrCode = 2
	Unknown            ErrCode = 3
	InvalidArgument    ErrCode = 4
	DeadlineExceeded   ErrCode = 5
	NotFound           ErrCode = 6
	AlreadyExists      ErrCode = 7
	PermissionDenied   ErrCode = 8
	ResourceExhausted  ErrCode = 9
	FailedPrecondition ErrCode = 10
	Aborted            ErrCode = 11
	OutOfRange         ErrCode = 12
	Unimplemented      ErrCode = 13
	Internal           ErrCode = 14
	Unavailable        ErrCode = 15
	DataLoss           ErrCode = 16
	Unauthenticated    ErrCode = 17
	TooManyRequests    ErrCode = 18
	InternalOnlyLog    ErrCode = 19
)

var codeNames = map[ErrCode]string{
	None:               "ok",
	NoContent:          "no_content",
	Canceled:           "canceled",
	Unknown:            "unknown",
	InvalidArgument:    "invalid_argument",
	DeadlineExceeded:   "deadline_exceeded",
	NotFound:           "not_found",
	AlreadyExists:      "already_exists",
	PermissionDenied:   "permission_denied",
	ResourceExhausted:  "resource_exhausted",
	FailedPrecondition: "failed_precondition",
	Aborted:            "aborted",
	OutOfRange:         "out_of_range",
	Unimplemented:      "unimplemented",
	Internal:           "internal",
	Unavailable:        "unavailable",
	DataLoss:           "data_loss",
	Unauthenticated:    "unauthenticated",
	TooManyRequests:    "too_many_requests",
	InternalOnlyLog:    "internal_only_log",
}

var httpStatusByCode = map[ErrCode]int{
	None:               http.StatusOK,
	NoContent:          http.StatusNoContent,
	Canceled:           http.StatusGatewayTimeout,
	Unknown:            http.StatusInternalServerError,
	InvalidArgument:    http.StatusBadRequest,
	DeadlineExceeded:   http.StatusGatewayTimeout,
	NotFound:           http.StatusNotFound,
	AlreadyExists:      http.StatusConflict,
	PermissionDenied:   http.StatusForbidden,
	ResourceExhausted:  http.StatusTooManyRequests,
	FailedPrecondition: http.StatusBadRequest,
	Aborted:            http.StatusConflict,
	OutOfRange:         http.StatusBadRequest,
	Unimplemented:      http.StatusNotImplemented,
	Internal:           http.StatusInternalServerError,
	Unavailable:        http.StatusServiceUnavailable,
	DataLoss:           http.StatusInternalServerError,
	Unauthenticated:    http.StatusUnauthorized,
	TooManyRequests:    http.StatusTooManyRequests,
	InternalOnlyLog:    http.StatusInternalServerError,
}

var grpcStatusByCode = map[ErrCode]codes.Code{
	None:               codes.OK,
	NoContent:          codes.OK,
	Canceled:           codes.Canceled,
	Unknown:            codes.Unknown,
	InvalidArgument:    codes.InvalidArgument,
	DeadlineExceeded:   codes.DeadlineExceeded,
	NotFound:           codes.NotFound,
	AlreadyExists:      codes.AlreadyExists,
	PermissionDenied:   codes.PermissionDenied,
	ResourceExhausted:  codes.ResourceExhausted,
	FailedPrecondition: codes.FailedPrecondition,
	Aborted:            codes.Aborted,
	OutOfRange:         codes.OutOfRange,
	Unimplemented:      codes.Unimplemented,
	Internal:           codes.Internal,
	Unavailable:        codes.Unavailable,
	DataLoss:           codes.DataLoss,
	Unauthenticated:    codes.Unauthenticated,
	TooManyRequests:    codes.ResourceExhausted,
	InternalOnlyLog:    codes.Internal,
}
