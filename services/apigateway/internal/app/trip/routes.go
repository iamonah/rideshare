package trip

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, handler *Handler) {
	if handler == nil {
		return
	}

	r.HandleFunc("/trip/preview", handler.HandlePreview).Methods(http.MethodPost)
}
