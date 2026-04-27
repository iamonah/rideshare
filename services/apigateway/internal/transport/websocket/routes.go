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

func RegisterRoutes(r *mux.Router, server *Server) {
	server.manager.RegisterHandler(driverEventTypes, server.ReceiveDriverEvents)
	server.manager.RegisterHandler(riderEventTypes, server.ReceiveRiderEvents)

	r.HandleFunc("/ws/drivers", server.HandleDriversWebsocket).Methods(http.MethodGet)
	r.HandleFunc("/ws/riders", server.HandleRidersWebsocket).Methods(http.MethodGet)
}
