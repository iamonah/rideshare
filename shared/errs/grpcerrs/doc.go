// Package grpcerrs adapts the core errs model to gRPC status errors.
//
// Mapping rules:
//   - AppError with fields -> code/message plus BadRequest details
//   - non-internal AppError without fields -> code/message as a normal gRPC status
//   - internal-class AppError without fields -> Internal with a safe internal service message
//   - plain error -> Internal with a safe internal service message
package grpcerrs
