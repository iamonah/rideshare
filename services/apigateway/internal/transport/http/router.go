package httptransport

import (
	"net/http"

	"github.com/gorilla/mux"
	driverapp "github.com/iamonah/rideshare/services/apigateway/internal/app/driver"
	paymentapp "github.com/iamonah/rideshare/services/apigateway/internal/app/payment"
	tripapp "github.com/iamonah/rideshare/services/apigateway/internal/app/trip"
	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	websockettransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/websocket"
)

type Dependencies struct {
	Trips      tripapp.PreviewTripUpstream
	Websockets *websockettransport.Handler
}

func NewRouter(deps Dependencies) http.Handler {
	router := mux.NewRouter()

	tripapp.RegisterRoutes(router, tripapp.NewHandler(deps.Trips))
	driverapp.RegisterRoutes(router, driverapp.NewHandler())
	paymentapp.RegisterRoutes(router, paymentapp.NewHandler())
	websockettransport.RegisterRoutes(router, deps.Websockets)

	return httpcommon.WithCORS(router)
}
