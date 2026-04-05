package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/iamonah/rideshare/shared/proto/pb/driverpb"
	"google.golang.org/grpc"
)

type DriverClient struct {
	client driverpb.DriverServiceClient
	conn   *grpc.ClientConn
}

func NewDriverClient(url string, opts ...grpc.DialOption) (*DriverClient, error) {
	if url == "" {
		return nil, errors.New("newDriverClient: url is required")
	}

	conn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("newDriverClient: %w", err)
	}

	return &DriverClient{
		client: driverpb.NewDriverServiceClient(conn),
		conn:   conn,
	}, nil
}

func (dc *DriverClient) RegisterDriver(ctx context.Context, req *driverpb.RegisterDriverRequest) (*driverpb.Driver, error) {
	if dc == nil || dc.client == nil {
		return nil, errors.New("registerDriver: client is not initialized")
	}

	resp, err := dc.client.RegisterDriver(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("registerDriver: %w", err)
	}

	return resp.GetDriver(), nil
}

func (dc *DriverClient) UnregisterDriver(ctx context.Context, req *driverpb.RegisterDriverRequest) (*driverpb.Driver, error) {
	if dc == nil || dc.client == nil {
		return nil, errors.New("unregisterDriver: client is not initialized")
	}

	resp, err := dc.client.UnregisterDriver(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("unregisterDriver: %w", err)
	}

	return resp.GetDriver(), nil
}

func (dc *DriverClient) Close() {
	if dc != nil && dc.conn != nil {
		if err := dc.conn.Close(); err != nil {
			fmt.Printf("closeDriverClient: %v\n", err)
		}
	}
}
