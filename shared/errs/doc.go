// Package errs defines the core error model for rideshare.
//
// The package uses one canonical error type, AppError, plus a FieldErrors
// builder for request validation details.
//
// AppError carries the client-facing code and message you want to return.
// FieldErrors can be attached for structured InvalidArgument responses.
//
// Transport adapters should preserve AppError codes/messages and fall back to a
// generic internal response for plain error values. See README.md in this
// directory.
package errs
