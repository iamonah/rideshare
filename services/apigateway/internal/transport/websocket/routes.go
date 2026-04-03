package websockettransport

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, handler *Handler) {
	r.HandleFunc("/ws/drivers", handler.HandleDriversWebsocket).Methods(http.MethodGet)
	r.HandleFunc("/ws/riders", handler.HandleRidersWebsocket).Methods(http.MethodGet)
}
