package trip

import (
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type PackageSlug string

var packageSlugs = make(map[string]PackageSlug)

func newPackageSlug(s string) PackageSlug {
	ps := PackageSlug(s)
	packageSlugs[strings.ToLower(s)] = ps
	return ps
}
func (p PackageSlug) String() string {
	return string(p)
}

var (
	PackageSlugVan    = newPackageSlug("van")
	PackageSlugSUV    = newPackageSlug("suv")
	PackageSlugSedan  = newPackageSlug("sedan")
	PackageSlugLuxury = newPackageSlug("luxury")
)

func ParsePackageSlug(s string) (PackageSlug, bool) {
	ps, ok := packageSlugs[strings.ToLower(s)]
	return ps, ok
}

// “The driver service is the source of truth for live driver data, but
// my trip service owns its own local representation of the driver data it needs.”
type AssignedDriverSnapshot struct {
	ID          string `bson:"id" json:"id"`
	FirstName   string `bson:"first_name" json:"first_name"`
	LastName    string `bson:"last_name" json:"last_name"`
	PhoneNumber string `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	VehicleID   string `bson:"vehicle_id,omitempty" json:"vehicle_id,omitempty"`
}

type Trip struct {
	ID       bson.ObjectID           `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string                  `bson:"user_id" json:"user_id"`
	Status   string                  `bson:"status" json:"status"`
	RideFare *RideFare               `bson:"ride_fare" json:"ride_fare"`
	Driver   *AssignedDriverSnapshot `bson:"driver" json:"driver"`
	// Route    *RouteSummary `bson:"route" json:"route"`
}
