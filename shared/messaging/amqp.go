package messaging

import "encoding/json"

// AmqpMessage is the message structure for AMQP.
type AmqpMessage struct {
	OwnerID string          `json:"ownerId"`
	Data    json.RawMessage `json:"data"`
}

const (
	// RideShareExchange is the single topic exchange used for all application messages.
	// Publishers only choose a routing key; queue bindings decide who receives it.
	RideShareExchange = "rideshare.messages"

	ErrorType = "error"  //for websocket error messages
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
	FindAvailableDriversQueue        = "find_available_drivers"         // driver-service consumes trip creation / retry events and chooses a driver.
	DriverCmdTripRequestQueue        = "driver_cmd_trip_request"        // api-gateway pushes trip requests to the connected driver websocket.
	DriverTripResponseQueue          = "driver_trip_response"           // trip workflow consumes driver accept/decline responses.
	NotifyRiderNoDriversFoundQueue   = "notify_rider_no_drivers_found"  // rider notification when no suitable driver was found.
	NotifyDriverAssignQueue          = "notify_driver_assign"           // rider notification when a driver has been assigned.
	PaymentTripResponseQueue         = "payment_trip_response"          // reserved for payment responses that affect trip state.
	NotifyPaymentSessionCreatedQueue = "notify_payment_session_created" // rider notification when a payment session is created.
	NotifyPaymentSuccessQueue        = "payment_success"                // rider notification when payment succeeds.
)
