package messaging

import "fmt"

func (rm *RabbitMQClient) setupSharedInfrastructure() error {
	if err := rm.declareExchange(RideShareExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("exchange %q: %w", RideShareExchange, err)
	}

	// Driver matching starts when a trip is created and can be retried when a driver rejects it.
	if err := rm.declareQueueAndBind(FindAvailableDriversQueue, []string{
		DriverEventDriverNotInterested,
		TripEventCreated,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", FindAvailableDriversQueue, err)
	}

	// Driver-facing trip request notifications are delivered by the websocket gateway.
	if err := rm.declareQueueAndBind(DriverCmdTripRequestQueue, []string{
		DriverCmdTripRequest,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", DriverCmdTripRequestQueue, err)
	}

	// Driver responses are consumed by the trip workflow so it can assign or retry.
	if err := rm.declareQueueAndBind(DriverTripResponseQueue, []string{
		DriverCmdTripAccept,
		DriverCmdTripDecline,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", DriverTripResponseQueue, err)
	}

	// Rider-facing notifications are split by intent so the websocket gateway can consume them independently.
	if err := rm.declareQueueAndBind(NotifyRiderNoDriversFoundQueue, []string{
		DriverEventNoDriversFound,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", NotifyRiderNoDriversFoundQueue, err)
	}

	if err := rm.declareQueueAndBind(NotifyDriverAssignQueue, []string{
		DriverEventDriverAssigned,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", NotifyDriverAssignQueue, err)
	}

	if err := rm.declareQueueAndBind(NotifyPaymentSessionCreatedQueue, []string{
		PaymentEventSessionCreated,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", NotifyPaymentSessionCreatedQueue, err)
	}

	if err := rm.declareQueueAndBind(NotifyPaymentSuccessQueue, []string{
		PaymentEventSuccess,
	}, nil); err != nil {
		return fmt.Errorf("bind queue %q : %w", NotifyPaymentSuccessQueue, err)
	}

	return nil
}
