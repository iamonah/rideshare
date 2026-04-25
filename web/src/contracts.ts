import { Coordinate, Driver, Route, RouteFare, Trip } from "./types";


// These are the endpoints the API Gateway must have for the frontend to work correctly
export enum BackendEndpoints {
  PREVIEW_TRIP = "/trip/preview",
  START_TRIP = "/trip/start",
  WS_DRIVERS = "/drivers",
  WS_RIDERS = "/riders",
}

export enum TripEvents {
  NoDriversFound = "driver.event.no_drivers_found",
  DriverAssigned = "driver.event.driver_assigned",
  Completed = "trip.event.completed",
  Cancelled = "trip.event.cancelled",
  Created = "trip.event.created",
  DriverLocation = "driver.cmd.location",
  DriverTripRequest = "driver.cmd.trip_request",
  DriverTripAccept = "driver.cmd.trip_accept",
  DriverTripDecline = "driver.cmd.trip_decline",
  DriverRegister = "driver.cmd.register",
  PaymentSessionCreated = "payment.event.session_created",
}

// Messages sent from the server to the client via the websocket
export type ServerWsMessage =
  | PaymentSessionCreatedRequest
  | DriverAssignedRequest
  | DriverLocationRequest
  | DriverTripRequest
  | DriverRegisterRequest
  | TripCreatedRequest
  | NoDriversFoundRequest;

// Messages sent from the client to the server via the websocket
export type ClientWsMessage = DriverResponseToTripResponse

interface TripCreatedRequest {
  type: TripEvents.Created;
  data: Trip;
}

interface NoDriversFoundRequest {
  type: TripEvents.NoDriversFound;
}

interface DriverRegisterRequest {
  type: TripEvents.DriverRegister;
  data: Driver | DriverWirePayload;
}
interface DriverTripRequest {
  type: TripEvents.DriverTripRequest;
  data: Trip | TripWirePayload;
}

export interface PaymentEventSessionCreatedData {
  tripID: string;
  sessionID: string;
  amount: number;
  currency: string;
}

interface PaymentSessionCreatedRequest {
  type: TripEvents.PaymentSessionCreated;
  data: PaymentEventSessionCreatedData;
}

interface DriverAssignedRequest {
  type: TripEvents.DriverAssigned;
  data: Trip | TripWirePayload;
}

interface DriverLocationRequest {
  type: TripEvents.DriverLocation;
  data: Driver[];
}

interface DriverResponseToTripResponse {
  type: TripEvents.DriverTripAccept | TripEvents.DriverTripDecline;
  data: {
    tripID: string;
    riderID: string;
    driver: Driver;
  };
}

interface DriverWirePayload {
  id: string;
  name: string;
  profilePicture?: string;
  profile_picture?: string;
  carPlate?: string;
  car_plate?: string;
  geohash?: string;
  geo_hash?: string;
  location?: Coordinate;
}

interface TripWirePayload {
  id?: string;
  tripId?: string;
  userID?: string;
  userId?: string;
  status: string;
  selectedFare?: RouteFare;
  fare?: {
    id: string;
    packageSlug: RouteFare["packageSlug"];
  };
  route: Route | {
    geometry: Coordinate[];
    duration?: number;
    distance?: number;
  };
  durationSeconds?: number;
  distanceMeters?: number;
  driver?: Driver | DriverWirePayload;
}

export interface HTTPTripPreviewResponse {
  route: Route;
  rideFares: RouteFare[];
}

export interface HTTPTripStartRequestPayload {
  rideFareId: string;
  userId: string;
}

export interface HTTPTripPreviewRequestPayload {
  userId: string;
  pickup: Coordinate;
  destination: Coordinate;
}

export function isValidTripEvent(event: string): event is TripEvents {
  return Object.values(TripEvents).includes(event as TripEvents);
}

export function isValidWsMessage(message: ServerWsMessage): message is ServerWsMessage {
  return isValidTripEvent(message.type);
}

export function normalizeDriver(data: Driver | DriverWirePayload): Driver {
  const wire = data as DriverWirePayload;

  return {
    id: data.id,
    name: data.name,
    profilePicture: wire.profilePicture ?? wire.profile_picture ?? "",
    carPlate: wire.carPlate ?? wire.car_plate ?? "",
    geohash: wire.geohash ?? wire.geo_hash ?? "",
    location: data.location ?? { latitude: 0, longitude: 0 },
  };
}

export function normalizeTrip(data: Trip | TripWirePayload): Trip {
  const wire = data as TripWirePayload;
  const route = wire.route;
  const geometry = Array.isArray(route.geometry)
    ? "coordinates" in (route.geometry[0] ?? {})
      ? route.geometry as { coordinates: Coordinate[] }[]
      : [{ coordinates: route.geometry as Coordinate[] }]
    : [];
  const duration = route.duration ?? wire.durationSeconds ?? 0;
  const distance = route.distance ?? wire.distanceMeters ?? 0;

  return {
    id: wire.id ?? wire.tripId ?? "",
    userID: wire.userID ?? wire.userId ?? "",
    status: wire.status,
    selectedFare: wire.selectedFare ?? {
      id: wire.fare?.id ?? "",
      packageSlug: wire.fare?.packageSlug ?? RouteFarePackageFallback,
      expiresAt: new Date(),
      route: {
        geometry,
        duration,
        distance,
      },
    },
    route: {
      geometry,
      duration,
      distance,
    },
    driver: wire.driver ? normalizeDriver(wire.driver) : undefined,
  };
}

const RouteFarePackageFallback = "sedan" as RouteFare["packageSlug"];
