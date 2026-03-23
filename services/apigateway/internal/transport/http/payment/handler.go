package paymenthttp

import apppayment "github.com/iamonah/rideshare/services/apigateway/internal/app/payment"

type Handler struct {
	payments *apppayment.Service
}

func NewHandler(payments *apppayment.Service) *Handler {
	return &Handler{payments: payments}
}
