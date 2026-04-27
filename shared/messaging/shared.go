package messaging

import "fmt"

func (rm *RabbitMQClient) setupSharedInfrastructure() error {
	// One shared topic exchange keeps publishers simple while routing remains expressive.
	if err := rm.declareExchange(RideShareExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup exchange %q: %w", RideShareExchange, err)
	}

	// Driver matching starts when a trip is created and can be retried when a driver rejects it.
	if err := rm.declareQueueAndBind(RideShareExchange, FindAvailableDriversQueue, []string{
		TripEventCreated,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", FindAvailableDriversQueue, RideShareExchange, err)
	}

	// Driver-facing trip request notifications are delivered by the websocket gateway.
	if err := rm.declareQueueAndBind(RideShareExchange, DriverCmdTripRequestQueue, []string{
		DriverCmdTripRequest,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", DriverCmdTripRequestQueue, RideShareExchange, err)
	}

	// Driver responses are consumed by the trip workflow so it can assign or retry.
	if err := rm.declareQueueAndBind(RideShareExchange, DriverTripResponseQueue, []string{
		DriverCmdTripAccept,
		DriverCmdTripDecline,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", DriverTripResponseQueue, RideShareExchange, err)
	}

	// A rejected driver feeds back into matching so a different driver can be selected.
	if err := rm.declareQueueAndBind(RideShareExchange, FindAvailableDriversQueue, []string{
		DriverEventDriverNotInterested,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", FindAvailableDriversQueue, RideShareExchange, err)
	}

	// Rider-facing notifications are split by intent so the websocket gateway can consume them independently.
	if err := rm.declareQueueAndBind(RideShareExchange, NotifyDriverNoDriversFoundQueue, []string{
		DriverEventNoDriversFound,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", NotifyDriverNoDriversFoundQueue, RideShareExchange, err)
	}

	if err := rm.declareQueueAndBind(RideShareExchange, NotifyDriverAssignQueue, []string{
		DriverEventDriverAssigned,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", NotifyDriverAssignQueue, RideShareExchange, err)
	}

	if err := rm.declareQueueAndBind(RideShareExchange, NotifyPaymentSessionCreatedQueue, []string{
		PaymentEventSessionCreated,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", NotifyPaymentSessionCreatedQueue, RideShareExchange, err)
	}

	if err := rm.declareQueueAndBind(RideShareExchange, NotifyPaymentSuccessQueue, []string{
		PaymentEventSuccess,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", NotifyPaymentSuccessQueue, RideShareExchange, err)
	}

	return nil
}
