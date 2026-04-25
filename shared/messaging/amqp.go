package messaging

import "encoding/json"

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string          `json:"ownerId"`
	Data    json.RawMessage `json:"data"`
}

// Routing keys - using consistent event/command patterns
const (
	// Exchanges are named by domain and message category.
	TripEventsExchange      = "trip.events"
	DriverCommandsExchange  = "driver.commands"
	DriverEventsExchange    = "driver.events"
	PaymentEventsExchange   = "payment.events"
	PaymentCommandsExchange = "payment.commands"

	// Trip events (trip.event.*)
	TripEventCreated   = "trip.event.created"
	TripEventUpdated   = "trip.event.updated"
	TripEventCancelled = "trip.event.cancelled"
	TripEventCompleted = "trip.event.completed"

	// Driver events (driver.event.*)
	DriverEventDriverAssigned      = "driver.event.driver_assigned"
	DriverEventNoDriversFound      = "driver.event.no_drivers_found"
	DriverEventDriverNotInterested = "driver.event.driver_not_interested"

	// Driver commands (driver.cmd.*)
	DriverCmdTripRequest = "driver.cmd.trip_request"
	DriverCmdTripAccept  = "driver.cmd.trip_accept"
	DriverCmdTripDecline = "driver.cmd.trip_decline"
	DriverCmdLocation    = "driver.cmd.location"
	DriverCmdRegister    = "driver.cmd.register"

	// Payment events (payment.event.*)
	PaymentEventSessionCreated = "payment.event.session_created"
	PaymentEventSuccess        = "payment.event.success"
	PaymentEventFailed         = "payment.event.failed"
	PaymentEventCancelled      = "payment.event.cancelled"

	// Payment commands (payment.cmd.*)
	PaymentCmdCreateSession = "payment.cmd.create_session"
)

const (
	// Queues
	DriverTripEventsQueue   = "driver.trip-events.queue"   // driver-service consumes trip/driver matching events.
	DriverTripRequestsQueue = "driver.trip-requests.queue" // api-gateway delivers trip requests to connected drivers.
	TripDriverEventsQueue   = "trip.driver-events.queue"   // trip-service consumes driver-side outcome events.
	TripDriverCommandsQueue = "trip.driver-commands.queue" // trip-service consumes driver response commands.
)
