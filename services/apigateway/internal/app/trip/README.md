# Trip Transport Naming

This document defines the client-facing JSON naming contract for the trip flow.

## Source Of Truth

The source of truth is the protobuf schema in
`shared/proto/tripservice.proto`.

- Proto field names stay in `snake_case`
- Client-facing JSON uses ProtoJSON `lowerCamelCase`
- Gateway HTTP DTOs must use the same `lowerCamelCase` names
- gRPC-side validation structs must use the same `lowerCamelCase` names
- Validation error field paths must use the same `lowerCamelCase` names

## Naming Rule

Use the ProtoJSON name derived from the proto field, not the raw proto field
name and not custom Go-style initialism names.

Examples:

- `user_id` -> `userId`
- `trip_id` -> `tripId`
- `ride_fare_id` -> `rideFareId`
- `ride_fares` -> `rideFares`
- `package_slug` -> `packageSlug`
- `total_price_in_cents` -> `totalPriceInCents`

Nested field paths must follow the same rule:

- `pickup.latitude`
- `pickup.longitude`
- `destination.latitude`
- `destination.longitude`

## What To Keep In Sync

These layers should all use the same field names:

1. HTTP request DTOs
2. HTTP response DTOs
3. Internal validation structs used before calling business logic
4. gRPC request-to-validation mapping
5. gRPC response-to-HTTP mapping
6. Validation error payloads
7. gRPC bad request field violations

## Example

Proto:

```proto
message PreviewTripRequest {
  string user_id = 1;
  Coordinate pickup = 2;
  Coordinate destination = 3;
}
```

Client-facing JSON:

```json
{
  "userId": "user-123",
  "pickup": {
    "latitude": 6.5244,
    "longitude": 3.3792
  },
  "destination": {
    "latitude": 6.6018,
    "longitude": 3.3515
  }
}
```

Go DTO:

```go
type PreviewTripInput struct {
	UserID      string           `json:"userId" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}
```

Validation errors should come back like:

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "validation failed",
    "fields": [
      {
        "field": "pickup.latitude",
        "message": "field is required"
      }
    ]
  }
}
```

## Important Note About Generated Protobuf Go Code

Do not use the generated `pb.go` struct `json` tags as the frontend contract.

Why:

- Generated protobuf Go structs often show `json:"snake_case"` tags
- ProtoJSON still serializes fields using `lowerCamelCase`
- Our HTTP contract should follow ProtoJSON naming, not `encoding/json` on the
  generated structs

In practice:

- use proto definitions as the naming source of truth
- use `lowerCamelCase` for gateway JSON tags
- use the same `lowerCamelCase` names in validation structs
- keep error field paths aligned with the same names

## Drift Warning

We currently duplicate some transport structs outside the proto-generated types.
That is fine for validation and transport boundaries, but those structs can drift
if the proto changes.

When you add or rename a proto field, update:

1. Gateway DTO `json` tags
2. gRPC validation structs
3. Mapping code
4. Validation/error assertions in tests
