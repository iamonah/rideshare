# trip service

This service handles all trip-related operations in the system.

## Architecture

The service follows a feature-first structure where the trip domain owns its
business types, ports, and use-case logic, while infra packages implement the
outer adapters.

```
services/trip-service/
├── cmd/                    # Application entry points
│   └── main.go            # Main application setup
├── internal/              # Private application code
│   ├── domain/
│   │   └── trip/         # Trip business types, ports, and use cases
│   └── infra/            # External adapter implementations
│       ├── events/       # Event handling
│       ├── external/     # Third-party clients like OSRM
│       ├── grpc/         # gRPC handlers and mappers
│       └── tripdb/       # Trip persistence
└── README.md            # This file
```

### Layer Responsibilities

1. **Domain Layer** (`internal/domain/trip/`)
   - Defines trip business types like routes, fares, and trips
   - Declares repository/provider contracts
   - Implements trip use cases in the same feature package

2. **Infrastructure Layer** (`internal/infra/`)
   - `tripdb/`: Implements data persistence
   - `events/`: Handles event publishing and consuming
   - `grpc/`: Handles gRPC communication
   - `external/`: Integrates external providers

## Key Benefits

1. **Dependency Inversion**: Outer adapters depend on trip domain contracts
2. **Separation of Concerns**: Business logic stays in the trip feature package
3. **Testability**: Easy to mock dependencies for testing
4. **Maintainability**: Clear boundaries between components
5. **Flexibility**: Easy to swap providers without affecting trip logic
