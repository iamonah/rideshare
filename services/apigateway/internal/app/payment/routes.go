package payment

import "github.com/gorilla/mux"

func RegisterRoutes(r *mux.Router, handler *Handler) {
	if r == nil || handler == nil {
		return
	}
}
