# Messaging Topology

This document explains the RabbitMQ naming rules used by the rideshare services.
The goal is to keep exchanges, routing keys, and queues easy to reason about as
the system grows.

## Mental Model

- Events are facts: something already happened.
- Commands are instructions: something should do work.
- Exchanges describe the producer domain and message category.
- Routing keys describe the specific event or command.
- Queues describe the consumer, audience, or processing purpose.

## Current Exchanges

```text
trip.events
driver.events
driver.commands
payment.events
payment.commands
```

Event exchanges are owned by the domain that discovered or owns the fact.
Command exchanges are owned by the target side responsible for carrying out the
action.

## Current Routing Keys

Trip events:

```text
trip.event.created
trip.event.updated
trip.event.cancelled
trip.event.completed
```

Driver events:

```text
driver.event.driver_assigned
driver.event.no_drivers_found
driver.event.driver_not_interested
```

Driver commands:

```text
driver.cmd.trip_request
driver.cmd.trip_accept
driver.cmd.trip_decline
driver.cmd.location
driver.cmd.register
```

Payment events and commands:

```text
payment.event.session_created
payment.event.success
payment.event.failed
payment.event.cancelled
payment.cmd.create_session
```

## Current Queues

```text
driver.trip-events.queue
driver.trip-requests.queue
trip.driver-events.queue
trip.driver-commands.queue
```

`driver.trip-events.queue` is consumed by driver-service. It carries trip events
that driver-service reacts to, such as `trip.event.created`.

`driver.trip-requests.queue` is a driver-facing delivery queue. The technical
consumer may be the API gateway websocket layer, but the business audience is the
driver.

`trip.driver-events.queue` is consumed by trip-service. It carries driver events
that trip-service reacts to, such as `driver.event.no_drivers_found`.

`trip.driver-commands.queue` is consumed by trip-service. It carries driver
commands that trip-service reacts to, such as `driver.cmd.trip_accept`.

## Current Bindings

```text
trip.events     + trip.event.created                 -> driver.trip-events.queue
driver.events   + driver.event.driver_not_interested -> driver.trip-events.queue
driver.commands + driver.cmd.trip_request            -> driver.trip-requests.queue
driver.events   + driver.event.no_drivers_found      -> trip.driver-events.queue
driver.events   + driver.event.driver_assigned       -> trip.driver-events.queue
driver.commands + driver.cmd.trip_accept             -> trip.driver-commands.queue
driver.commands + driver.cmd.trip_decline            -> trip.driver-commands.queue
```

## Naming Rules

Routing keys use this pattern:

```text
<producer-domain>.<message-kind>.<message-name>
```

Examples:

```text
trip.event.created
driver.event.no_drivers_found
driver.cmd.trip_request
```

Service event queues use this pattern:

```text
<consumer-domain>.<source-domain>-events.queue
```

Examples:

```text
driver.trip-events.queue
trip.driver-events.queue
```

Audience or workflow queues use this pattern:

```text
<audience>.<purpose>.queue
```

Example:

```text
driver.trip-requests.queue
```

## Event Ownership

Name events by the service that owns or discovered the fact, not merely by what
the event is about.

For example, `no_drivers_found` is a driver-service matching result. It should be:

```text
driver.event.no_drivers_found
```

not:

```text
trip.event.no_drivers_found
```

Trip-service can still consume the driver event and update trip state. The event
name should represent the producer of the fact.

## Commands vs Events

Sending a trip offer to a driver is a command:

```text
driver.cmd.trip_request
```

It means: ask this driver to respond to the trip request.

The driver's response becomes an event:

```text
driver.event.driver_assigned
driver.event.driver_not_interested
```

That means the driver side has produced a new fact that other services can react
to.

## Message Contracts

Broker payloads should use small shared contract types, not full domain models
or gRPC transport types.

Event payloads should include only the fields consumers need. This keeps broker
contracts stable and avoids leaking internal service models across boundaries.
