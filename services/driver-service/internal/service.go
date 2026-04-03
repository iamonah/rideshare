package driverservice

import (
	"errors"
	"strings"
	"sync"

	"github.com/iamonah/rideshare/shared/errs"
	driverpb "github.com/iamonah/rideshare/shared/proto/pb/driverpb"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	mu      sync.RWMutex
	drivers map[string]*driverpb.Driver
}

func NewService() *Service {
	return &Service{
		drivers: make(map[string]*driverpb.Driver),
	}
}

func (s *Service) RegisterDriver(req *driverpb.RegisterDriverRequest) (*driverpb.Driver, error) {
	driverID := strings.TrimSpace(req.GetDriverId())
	packageSlug := strings.TrimSpace(req.GetPackageSlug())

	driver := &driverpb.Driver{
		Id:          driverID,
		PackageSlug: packageSlug,
	}

	s.mu.Lock()
	s.drivers[driverID] = driver
	s.mu.Unlock()

	return cloneDriver(driver), nil
}

func (s *Service) UnregisterDriver(req *driverpb.RegisterDriverRequest) (*driverpb.Driver, error) {
	driverID := strings.TrimSpace(req.GetDriverId())
	packageSlug := strings.TrimSpace(req.GetPackageSlug())

	s.mu.Lock()
	defer s.mu.Unlock()

	driver, ok := s.drivers[driverID]
	if !ok || driver.GetPackageSlug() != packageSlug {
		return nil, errs.New(errs.NotFound, errors.New("driver not found"))
	}

	delete(s.drivers, driverID)
	return cloneDriver(driver), nil
}

func cloneDriver(driver *driverpb.Driver) *driverpb.Driver {
	if driver == nil {
		return nil
	}

	return proto.Clone(driver).(*driverpb.Driver)
}
