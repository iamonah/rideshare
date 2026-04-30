package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamonah/rideshare/services/payment-service/internal/events"
	"github.com/iamonah/rideshare/services/payment-service/internal/infrastructure/stripe"
	"github.com/iamonah/rideshare/services/payment-service/internal/service"
	"github.com/iamonah/rideshare/services/payment-service/pkg/types"
	"github.com/iamonah/rideshare/shared/env"
	"github.com/iamonah/rideshare/shared/messaging"
)

var (
	rabbitUsername = env.GetString("RABBITMQ_DEFAULT_USER", "")
	rabbitPassword = env.GetString("RABBITMQ_DEFAULT_PASS", "")
	rabbitHost     = env.GetString("RABBITMQ_HOST", "")
	rabbitVhost    = env.GetString("RABBITMQ_VHOST", "")
	rabbitPort     = env.GetInt("RABBITMQ_PORT", 5672)
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	appURL := env.GetString("APP_URL", "http://localhost:3000")

	// Stripe config
	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: env.GetString("STRIPE_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel"),
	}

	if stripeCfg.StripeSecretKey == "" {
		exitWithError("STRIPE_SECRET_KEY is not set")
	}

	// Stripe processor
	paymentProcessor := stripe.NewStripeClient(stripeCfg)

	// Service
	svc := service.NewPaymentService(paymentProcessor)

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQClient(messaging.RabbitMqConfig{
		Username: rabbitUsername,
		Password: rabbitPassword,
		Host:     rabbitHost,
		Vhost:    rabbitVhost,
		Port:     int16(rabbitPort),
	})
	if err != nil {
		exitWithError("failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.Close()

	// Trip Consumer
	tripConsumer := events.NewTripConsumer(rabbitmq, svc)
	if err := tripConsumer.Listen(ctx); err != nil {
		exitWithError("failed to register payment consumer: %v", err)
	}

	// Wait for shutdown signal
	<-ctx.Done()
}

func exitWithError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
