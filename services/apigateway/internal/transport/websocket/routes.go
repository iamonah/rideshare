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

func RegisterRoutes(r *mux.Router, handler *Handler) {
	handler.Manager.RegisterHandler(driverEventTypes, handler.ReceiveDriverEvents)

	r.HandleFunc("/ws/drivers", handler.HandleDriversWebsocket).Methods(http.MethodGet)
	r.HandleFunc("/ws/riders", handler.HandleRidersWebsocket).Methods(http.MethodGet)
}
