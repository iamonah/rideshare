package driverhttp

import appdriver "github.com/iamonah/rideshare/services/apigateway/internal/app/driver"

type Handler struct {
	drivers *appdriver.Service
}

func NewHandler(drivers *appdriver.Service) *Handler {
	return &Handler{drivers: drivers}
}
