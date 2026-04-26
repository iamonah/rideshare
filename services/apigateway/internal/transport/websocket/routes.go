package websockettransport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamonah/rideshare/shared/messaging"
)

var driverEventTypes = []string{
	messaging.DriverCmdLocation,
	messaging.DriverCmdTripAccept,
	messaging.DriverCmdTripDecline,
}

var riderEventTypes = []string{
	// Define rider event types in the future
}

func RegisterRoutes(r *mux.Router, handler *Handler) {
	handler.Manager.RegisterHandler(driverEventTypes, handler.ReceiveDriverEvents)
	handler.Manager.RegisterHandler(riderEventTypes, handler.ReceiveRiderEvents)

	r.HandleFunc("/ws/drivers", handler.HandleDriversWebsocket).Methods(http.MethodGet)
	r.HandleFunc("/ws/riders", handler.HandleRidersWebsocket).Methods(http.MethodGet)
}
