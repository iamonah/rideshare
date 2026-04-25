package messaging

import "fmt"

func (rm *RabbitMQClient) setupSharedInfrastructure() error {
	// Exchanges
	if err := rm.declareExchange(TripEventsExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup exchange %q: %w", TripEventsExchange, err)
	}

	if err := rm.declareExchange(DriverCommandsExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup exchange %q: %w", DriverCommandsExchange, err)
	}

	if err := rm.declareExchange(DriverEventsExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup exchange %q: %w", DriverEventsExchange, err)
	}

	if err := rm.declareExchange(PaymentEventsExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup exchange %q: %w", PaymentEventsExchange, err)
	}

	if err := rm.declareExchange(PaymentCommandsExchange, TopicExchangeKind); err != nil {
		return fmt.Errorf("setup exchange %q: %w", PaymentCommandsExchange, err)
	}

	// Trip events
	if err := rm.declareQueueAndBind(TripEventsExchange, DriverTripEventsQueue, []string{
		TripEventCreated,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", DriverTripEventsQueue, TripEventsExchange, err)
	}

	// Driver commands
	if err := rm.declareQueueAndBind(DriverCommandsExchange, DriverTripRequestsQueue, []string{
		DriverCmdTripRequest,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", DriverTripRequestsQueue, DriverCommandsExchange, err)
	}

	if err := rm.declareQueueAndBind(DriverCommandsExchange, TripDriverCommandsQueue, []string{
		DriverCmdTripAccept,
		DriverCmdTripDecline,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", TripDriverCommandsQueue, DriverCommandsExchange, err)
	}

	// Driver events
	if err := rm.declareQueueAndBind(DriverEventsExchange, DriverTripEventsQueue, []string{
		DriverEventDriverNotInterested,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", DriverTripEventsQueue, DriverEventsExchange, err)
	}

	if err := rm.declareQueueAndBind(DriverEventsExchange, TripDriverEventsQueue, []string{
		DriverEventNoDriversFound,
		DriverEventDriverAssigned,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", TripDriverEventsQueue, DriverEventsExchange, err)
	}

	if err := rm.declareQueueAndBind(DriverEventsExchange, RiderEventsQueue, []string{
		DriverEventNoDriversFound,
		DriverEventDriverAssigned,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", RiderEventsQueue, DriverEventsExchange, err)
	}

	// Payment events
	if err := rm.declareQueueAndBind(PaymentEventsExchange, RiderEventsQueue, []string{
		PaymentEventSessionCreated,
	}, nil); err != nil {
		return fmt.Errorf("setup queue %q on %q: %w", RiderEventsQueue, PaymentEventsExchange, err)
	}

	return nil
}
