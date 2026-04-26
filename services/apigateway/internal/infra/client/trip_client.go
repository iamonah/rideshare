package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/iamonah/rideshare/services/apigateway/internal/app"
	"github.com/iamonah/rideshare/shared/proto/pb/trippb"
	"github.com/iamonah/rideshare/shared/types"
	"google.golang.org/grpc"
)

type Client struct {
	client trippb.TripServiceClient
	conn   *grpc.ClientConn
}

func NewClient(url string, opts ...grpc.DialOption) (*Client, error) {
	if url == "" {
		return nil, errors.New("newTripClient: url is required")
	}

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("newTripClient: %w", err)
	}

	return &Client{
		client: trippb.NewTripServiceClient(conn),
		conn:   conn,
	}, nil
}

func (c *Client) PreviewTrip(ctx context.Context, input app.PreviewTripInput) (*app.PreviewTripOutput, error) {
	req := toPreviewTripProto(input)

	resp, err := c.client.PreviewTrip(ctx, req)
	if err != nil {
		log.Printf("trip gRPC PreviewTrip failed: %v", err)
		return nil, err
	}

	return toPreviewTripOutput(resp), nil
}

func (c *Client) CreateTrip(ctx context.Context, input app.CreateTripInput) (*app.CreateTripOutput, error) {
	req := toCreateTripProto(input)

	resp, err := c.client.CreateTrip(ctx, req)
	if err != nil {
		log.Printf("trip gRPC CreateTrip failed: %v", err)
		return nil, err
	}

	return toCreateTripOutput(resp), nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("closeTripClient: %w", err)
	}
	return nil
}

func toCreateTripProto(input app.CreateTripInput) *trippb.CreateTripRequest {
	return &trippb.CreateTripRequest{
		RideFareId: input.RideFareID,
		UserId:     input.UserID,
	}
}

func toCreateTripOutput(resp *trippb.CreateTripResponse) *app.CreateTripOutput {
	if resp == nil {
		return &app.CreateTripOutput{}
	}

	return &app.CreateTripOutput{
		TripID: resp.GetTripId(),
		Trip:   toTrip(resp.GetTrip()),
	}
}

func toPreviewTripProto(input app.PreviewTripInput) *trippb.PreviewTripRequest {
	return &trippb.PreviewTripRequest{
		UserId: input.UserID,
		Pickup: &trippb.Coordinate{
			Latitude:  input.Pickup.Latitude,
			Longitude: input.Pickup.Longitude,
		},
		Destination: &trippb.Coordinate{
			Latitude:  input.Destination.Latitude,
			Longitude: input.Destination.Longitude,
		},
	}
}

func toPreviewTripOutput(resp *trippb.PreviewTripResponse) *app.PreviewTripOutput {
	if resp == nil {
		return &app.PreviewTripOutput{}
	}

	output := &app.PreviewTripOutput{
		TripID:    resp.GetTripId(),
		Route:     toRoute(resp.GetRoute()),
		RideFares: make([]app.PreviewRideFare, 0, len(resp.GetRideFares())),
	}

	for _, fare := range resp.GetRideFares() {
		output.RideFares = append(output.RideFares, app.PreviewRideFare{
			ID:                fare.GetId(),
			UserID:            fare.GetUserId(),
			PackageSlug:       fare.GetPackageSlug(),
			TotalPriceInCents: fare.GetTotalPriceInCents(),
			Route:             toRoute(fare.GetRoute()),
		})
	}

	return output
}

func toTrip(trip *trippb.Trip) *app.Trip {
	if trip == nil {
		return nil
	}

	return &app.Trip{
		ID:           trip.GetId(),
		SelectedFare: toRideFare(trip.GetSelectedFare()),
		Route:        toRoute(trip.GetRoute()),
		Status:       trip.GetStatus(),
		UserID:       trip.GetUserId(),
		// Driver:       toTripDriver(trip.GetDriver()),
	}
}

func toRideFare(fare *trippb.RideFare) *app.RideFare {
	if fare == nil {
		return nil
	}

	return &app.RideFare{
		ID:                fare.GetId(),
		UserID:            fare.GetUserId(),
		PackageSlug:       fare.GetPackageSlug(),
		TotalPriceInCents: fare.GetTotalPriceInCents(),
		Route:             toRoute(fare.GetRoute()),
	}
}

func toTripDriver(driver *trippb.TripDriver) *app.TripDriver {
	if driver == nil {
		return nil
	}

	return &app.TripDriver{
		ID:             driver.GetId(),
		Name:           driver.GetName(),
		ProfilePicture: driver.GetProfilePicture(),
		CarPlate:       driver.GetCarPlate(),
	}
}

func toRoute(route *trippb.Route) app.Route {
	if route == nil {
		return app.Route{}
	}

	geometry := make([]*app.Geometry, 0, len(route.GetGeometry()))
	for _, segment := range route.GetGeometry() {
		coordinates := make([]*types.Coordinate, 0, len(segment.GetCoordinates()))
		for _, coordinate := range segment.GetCoordinates() {
			coordinates = append(coordinates, &types.Coordinate{
				Latitude:  coordinate.GetLatitude(),
				Longitude: coordinate.GetLongitude(),
			})
		}

		geometry = append(geometry, &app.Geometry{
			Coordinates: coordinates,
		})
	}

	return app.Route{
		Distance: route.GetDistance(),
		Duration: route.GetDuration(),
		Geometry: geometry,
	}
}
