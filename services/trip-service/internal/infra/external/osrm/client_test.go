package osrm

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/types"
)

func TestGetRouteReturnsDeadlineExceededOnHTTPTimeout(t *testing.T) {
	t.Parallel()

	client := NewClient(&http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return nil, context.DeadlineExceeded
		}),
	}, "http://example.com")

	_, err := client.GetRoute(context.Background(), samplePickup(), sampleDestination())
	if err == nil {
		t.Fatal("expected timeout error")
	}

	var appErr *errs.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != errs.DeadlineExceeded {
		t.Fatalf("expected %v, got %v", errs.DeadlineExceeded, appErr.Code)
	}
}

func TestGetRouteReturnsCanceledWhenContextCanceled(t *testing.T) {
	t.Parallel()

	client := NewClient(&http.Client{
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, req.Context().Err()
		}),
	}, "http://example.com")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetRoute(ctx, samplePickup(), sampleDestination())
	if err == nil {
		t.Fatal("expected canceled error")
	}

	var appErr *errs.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != errs.Canceled {
		t.Fatalf("expected %v, got %v", errs.Canceled, appErr.Code)
	}
}

func samplePickup() *types.Coordinate {
	return &types.Coordinate{Latitude: 6.5244, Longitude: 3.3792}
}

func sampleDestination() *types.Coordinate {
	return &types.Coordinate{Latitude: 6.6018, Longitude: 3.3515}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
