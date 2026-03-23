package osrm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	tripdomain "github.com/iamonah/rideshare/services/trip-service/internal/domain/trip"
	"github.com/iamonah/rideshare/shared/errs"
	"github.com/iamonah/rideshare/shared/types"
)

const defaultBaseURL = "http://router.project-osrm.org"

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    httpClient,
	}
}

func (c *Client) GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*tripdomain.Route, error) {
	url := fmt.Sprintf(
		"%s/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		c.baseURL,
		pickup.Longitude, pickup.Latitude,
		destination.Longitude, destination.Latitude,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to build route request")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, errs.Newf(errs.Unavailable, err, "route provider unavailable")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to read route response")
	}

	if resp.StatusCode != http.StatusOK {
		var routeErr errorResponse
		if err := json.Unmarshal(body, &routeErr); err != nil {
			return nil, errs.Newf(
				errs.Internal,
				err,
				"json.Unmarshal route provider error response (status=%d)",
				resp.StatusCode,
			)
		}

		return nil, classifyError(resp.StatusCode, &routeErr)
	}

	var routeResp routeResponse
	if err := json.Unmarshal(body, &routeResp); err != nil {
		return nil, errs.Newf(errs.Internal, err, "failed to parse route response")
	}

	return toDomainRoute(routeResp), nil
}

func classifyError(statusCode int, routeErr *errorResponse) error {
	code := ""
	if routeErr != nil {
		code = routeErr.Code
	}

	switch code {
	case "NoRoute":
		return errs.New(errs.NotFound, errors.New("no route found between pickup and destination"))
	case "NoSegment":
		return errs.New(errs.InvalidArgument, errors.New("pickup or destination is not on a routable road"))
	case "InvalidValue", "InvalidQuery", "InvalidOptions", "InvalidUrl":
		return errs.New(errs.InvalidArgument, errors.New("route request is invalid"))
	}

	if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
		return errs.New(errs.InvalidArgument, errors.New("route request is invalid"))
	}

	return errs.Newf(
		errs.Unavailable,
		fmt.Errorf("route provider returned status %d with code %q", statusCode, code),
		"route provider unavailable",
	)
}

func toDomainRoute(routeResp routeResponse) *tripdomain.Route {
	routes := make([]tripdomain.RouteSummary, 0, len(routeResp.Routes))
	for _, route := range routeResp.Routes {
		routes = append(routes, tripdomain.RouteSummary{
			Distance: route.Distance,
			Duration: route.Duration,
			Geometry: tripdomain.RouteGeometry{
				Coordinates: route.Geometry.Coordinates,
			},
		})
	}

	return &tripdomain.Route{Routes: routes}
}
