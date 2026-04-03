package driverservice

import (
	"context"
	"errors"

	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/errs/grpcerrs"
	driverpb "github.com/iamonah/rideshare/shared/proto/pb/driverpb"
	"google.golang.org/grpc"
)

type GRPCHandler struct {
	driverpb.UnimplementedDriverServiceServer

	service *Service
}

type registerDriverInput struct {
	DriverID    string `json:"driver_id" validate:"required"`
	PackageSlug string `json:"package_slug" validate:"required"`
}

func NewGRPCHandler(s *grpc.Server, service *Service) {
	handler := &GRPCHandler{
		service: service,
	}

	driverpb.RegisterDriverServiceServer(s, handler)
}

func (h *GRPCHandler) RegisterDriver(ctx context.Context, req *driverpb.RegisterDriverRequest) (*driverpb.RegisterDriverResponse, error) {
	_ = ctx

	if err := validateRegisterDriverRequest(req); err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	driver, err := h.service.RegisterDriver(req)
	if err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	return &driverpb.RegisterDriverResponse{Driver: driver}, nil
}

func (h *GRPCHandler) UnregisterDriver(ctx context.Context, req *driverpb.RegisterDriverRequest) (*driverpb.RegisterDriverResponse, error) {
	_ = ctx

	if err := validateRegisterDriverRequest(req); err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	driver, err := h.service.UnregisterDriver(req)
	if err != nil {
		return nil, grpcerrs.ToStatus(err)
	}

	return &driverpb.RegisterDriverResponse{Driver: driver}, nil
}

func validateRegisterDriverRequest(req *driverpb.RegisterDriverRequest) error {
	if req == nil {
		return errs.New(errs.InvalidArgument, errors.New("request is required"))
	}

	input := registerDriverInput{
		DriverID:    req.GetDriverId(),
		PackageSlug: req.GetPackageSlug(),
	}

	if err := errs.Validate(input); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	return nil
}
