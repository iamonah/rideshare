package trip

type Route struct {
	Routes []RouteSummary
}

type RouteSummary struct {
	Distance float64
	Duration float64
	Geometry RouteGeometry
}

type RouteGeometry struct {
	Coordinates [][]float64
}
