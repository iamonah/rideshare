package httptransport

import (
	"net/http"

	"github.com/gorilla/mux"
	appdriver "github.com/iamonah/rideshare/services/apigateway/internal/app/driver"
	apppayment "github.com/iamonah/rideshare/services/apigateway/internal/app/payment"
	apptrip "github.com/iamonah/rideshare/services/apigateway/internal/app/trip"
	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	driverhttp "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/driver"
	paymenthttp "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/payment"
	triphttp "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/trip"
	websockettransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/websocket"
)

type Dependencies struct {
	Trips      *apptrip.Service
	Drivers    *appdriver.Service
	Payments   *apppayment.Service
	Websockets *websockettransport.Handler
}

func NewRouter(deps Dependencies) http.Handler {
	router := mux.NewRouter()

	triphttp.RegisterRoutes(router, triphttp.NewHandler(deps.Trips))
	driverhttp.RegisterRoutes(router, driverhttp.NewHandler(deps.Drivers))
	paymenthttp.RegisterRoutes(router, paymenthttp.NewHandler(deps.Payments))
	websockettransport.RegisterRoutes(router, deps.Websockets)

	return httpcommon.WithCORS(router)
}
