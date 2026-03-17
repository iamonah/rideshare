# `shared/errs`

`shared/errs` gives the codebase one canonical application error shape:
[`AppError`](/home/leo/Documents/rideshare/shared/errs/errs.go).

Callers should not construct `AppError` directly. Use the exported helpers that
set the exact code and message you want clients to receive.

## Rules

Use [`FieldErrors`](/home/leo/Documents/rideshare/shared/errs/field.go) when
you want structured validation details.

Use `New` when you want to promote an existing error into an `AppError`.

Use `Newf` when you want to keep an underlying error and add a safe formatted
message for the client.

At service boundaries, prefer returning `AppError` values instead of plain
`error` so the service decides the code, client message, and captured callsite.

Plain `error` is still supported as a fallback. Transport adapters will turn it
into a generic internal response.

## What Each Path Means

`AppError`:
- Carries the exact `ErrCode` and message to return to the client.
- May also include validation fields and captured callsite metadata.

Plain `error`:
- Is a fallback for unexpected failures that have not been promoted into an
  `AppError`.
- Is translated by transport adapters into a safe generic response.

## Constructors

`FieldErrors`

```go
fieldErrs := errs.NewFieldErrors()
fieldErrs.AddMessage("pickup", "is required")
fieldErrs.AddMessage("destination", "is required")

if err := fieldErrs.ToError(); err != nil {
	return err
}
```

`New`

```go
if trip == nil {
	return errs.New(errs.NotFound, errors.New("trip not found"))
}
```

`Newf`

```go
if err := repo.Create(ctx, trip); err != nil {
	return errs.Newf(errs.Internal, err, "failed to create trip")
}
```

`Validation`

```go
fieldErrs := errs.NewFieldErrors()
fieldErrs.AddMessage("request", "is required")
return errs.Validation(fieldErrs)
```

Plain internal error

```go
routeResp, err := client.GetRoute(ctx, pickup, destination)
if err != nil {
	return errs.Newf(errs.Unavailable, err, "route provider unavailable")
}
```

## gRPC Mapping

The gRPC adapter lives in
[`grpcerrs/status.go`](/home/leo/Documents/rideshare/shared/errs/grpcerrs/status.go).

It interprets errors like this:

- `AppError` with fields -> `ErrCode.GRPCStatus()` plus `BadRequest` details
- non-internal `AppError` without fields -> `ErrCode.GRPCStatus()` with the
  exact message
- internal-class `AppError` without fields -> `Internal` with the safe message
  `"internal service error"`
- plain `error` -> `Internal` with the safe message `"internal service error"`

Handler example:

```go
resp, err := svc.GetRoute(ctx, start, end)
if err != nil {
	return nil, grpcerrs.ToStatus(err)
}
```

## Logging

`AppError` captures `FuncName` and `FileName` for values created through the
package constructors.

Validation helpers also preserve the caller above the helper itself, so
`errs.Validation(...)` and `FieldErrors.ToError()` point back to the service or
handler line that created the error rather than to the shared `errs` package.

That metadata is for diagnostics only. It should be logged, not exposed to
clients or used for business logic.

Plain internal `error` values do not carry caller metadata. If you need origin
information for them, log at the failure site or introduce an explicit internal
error constructor later.
