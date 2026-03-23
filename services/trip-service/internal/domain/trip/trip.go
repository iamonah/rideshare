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

type Trip struct {
	ID       bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID   string        `bson:"user_id" json:"user_id"`
	Status   string        `bson:"status" json:"status"`
	RideFare *RideFare     `bson:"ride_fare" json:"ride_fare"`
}
