package websockettransport

import "github.com/iamonah/rideshare/services/apigateway/internal/infra/client"

type Handler struct {
	Manager      *Manager
	driverClient *client.DriverClient
}

func NewHandler(dc *client.DriverClient) *Handler {
	return &Handler{
		Manager:      NewManager(),
		driverClient: dc,
	}
}
