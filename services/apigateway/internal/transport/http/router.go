package httptransport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iamonah/rideshare/services/apigateway/internal/app"
	httpcommon "github.com/iamonah/rideshare/services/apigateway/internal/transport/http/common"
	websockettransport "github.com/iamonah/rideshare/services/apigateway/internal/transport/websocket"
	"github.com/iamonah/rideshare/shared/messaging"
)

type Dependencies struct {
	Handlers   app.TripUpstream
	Websockets *websockettransport.Server
	rabbitmq   *messaging.RabbitMQClient
}

func NewRouter(deps Dependencies) http.Handler {
	router := mux.NewRouter()

	app.RegisterRoutes(router, app.NewHandler(deps.Handlers))
	websockettransport.RegisterRoutes(router, deps.Websockets)

	return httpcommon.WithCORS(router)
}
