package grpc_client

import (
	"errors"
	"fmt"

	"github.com/iamonah/rideshare/shared/proto/pb/trip"
	"google.golang.org/grpc"
)

type TripClient struct {
	Client trip.TripServiceClient
	conn   *grpc.ClientConn
}

func NewTripClient(url string, opts ...grpc.DialOption) (*TripClient, error) {
	if url == "" {
		return nil, errors.New("newTripClient: url is required")
	}

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("newTripClient: %w", err)
	}

	c := trip.NewTripServiceClient(conn)

	return &TripClient{
		Client: c,
		conn:   conn,
	}, nil
}

func (c *TripClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("closeTripClient: %w", err)
	}
	return nil
}
