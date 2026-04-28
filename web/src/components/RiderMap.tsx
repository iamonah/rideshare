'use client';

import { useRiderStreamConnection } from '../hooks/useRiderStreamConnection';
import { MapContainer, Marker, Popup, Rectangle, TileLayer } from 'react-leaflet'
import L from 'leaflet';
import { getGeohashBounds } from '../utils/geohash';
import { useEffect, useMemo, useRef, useState } from 'react';
import { MapClickHandler } from './MapClickHandler';
import { Button } from './ui/button';
import { RouteFare, RequestRideProps, TripPreview, HTTPTripStartResponse } from "../types";
import { RoutingControl } from "./RoutingControl";
import { API_URL } from '../constants';
import { RiderTripOverview } from './RiderTripOverview';
import { BackendEndpoints, HTTPTripPreviewRequestPayload, HTTPTripPreviewResponse, HTTPTripStartRequestPayload } from '../contracts';
import { DriverAvatar } from './DriverAvatar';

const userMarker = new L.Icon({
    iconUrl: "/markers/location-pin.svg",
    iconSize: [40, 40],
    iconAnchor: [20, 40],
    popupAnchor: [0, -36],
});

const driverMarker = new L.Icon({
    iconUrl: "https://www.svgrepo.com/show/25407/car.svg",
    iconSize: [30, 30],
    iconAnchor: [15, 30],
});

interface RiderMapProps {
    onRouteSelected?: (distance: number) => void;
}

const TEST_RIDER_LOCATION = {
    latitude: 37.7749,
    longitude: -122.4194,
}

export default function RiderMap({ onRouteSelected }: RiderMapProps) {
    const [trip, setTrip] = useState<TripPreview | null>(null)
    const [previewError, setPreviewError] = useState<string | null>(null)
    const [selectedCarPackage] = useState<RouteFare | null>(null)
    const [destination, setDestination] = useState<[number, number] | null>(null)
    const mapRef = useRef<L.Map>(null)
    const userID = useMemo(() => crypto.randomUUID(), [])
    const debounceTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const previewAbortControllerRef = useRef<AbortController | null>(null);

    const {
        drivers,
        error,
        tripStatus,
        assignedDriver,
        paymentSession,
        resetTripStatus
    } = useRiderStreamConnection(TEST_RIDER_LOCATION, userID);

    console.log(tripStatus)

    useEffect(() => {
        if (!mapRef.current) {
            return
        }

        mapRef.current.setView([TEST_RIDER_LOCATION.latitude, TEST_RIDER_LOCATION.longitude], 13)
    }, [])

    const handleMapClick = async (e: L.LeafletMouseEvent) => {
        if (trip?.tripId) {
            return
        }

        if (debounceTimeoutRef.current) {
            clearTimeout(debounceTimeoutRef.current);
        }

        if (previewAbortControllerRef.current) {
            previewAbortControllerRef.current.abort()
            previewAbortControllerRef.current = null
        }

        debounceTimeoutRef.current = setTimeout(async () => {
            const abortController = new AbortController()
            previewAbortControllerRef.current = abortController
            setDestination([e.latlng.lat, e.latlng.lng])
            setPreviewError(null)

            try {
                const data = await requestRidePreview({
                    pickup: [TEST_RIDER_LOCATION.latitude, TEST_RIDER_LOCATION.longitude],
                    destination: [e.latlng.lat, e.latlng.lng],
                }, abortController.signal)

                if (previewAbortControllerRef.current !== abortController) {
                    return
                }

                const routeCoordinates = data.route?.geometry?.[0]?.coordinates
                if (!routeCoordinates || routeCoordinates.length === 0) {
                    throw new Error("Route preview is temporarily unavailable")
                }

                const parsedRoute = routeCoordinates
                    .map((coord) => [coord.latitude, coord.longitude] as [number, number])

                setTrip({
                    tripId: "",
                    route: parsedRoute,
                    rideFares: data.rideFares,
                    distance: data.route.distance,
                    duration: data.route.duration,
                })

                onRouteSelected?.(data.route.distance)
            } catch (err) {
                if (err instanceof DOMException && err.name === "AbortError") {
                    return
                }

                const message = err instanceof Error ? err.message : "Unable to preview this trip right now"
                setTrip(null)
                setPreviewError(message)
            } finally {
                if (previewAbortControllerRef.current === abortController) {
                    previewAbortControllerRef.current = null
                }
            }
        }, 500);
    }

    const requestRidePreview = async (props: RequestRideProps, signal?: AbortSignal): Promise<HTTPTripPreviewResponse> => {
        const { pickup, destination } = props
        const payload = {
            userId: userID,
            pickup: {
                latitude: pickup[0],
                longitude: pickup[1],
            },
            destination: {
                latitude: destination[0],
                longitude: destination[1],
            },
        } as HTTPTripPreviewRequestPayload

        const response = await fetch(`${API_URL}${BackendEndpoints.PREVIEW_TRIP}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            signal,
            body: JSON.stringify(payload),
        })

        const body = await response.json() as {
            data?: HTTPTripPreviewResponse
            error?: { message?: string }
        }
        if (!response.ok || !body.data) {
            throw new Error(body.error?.message || "Unable to preview this trip right now")
        }

        const { data } = body
        return data
    }

    useEffect(() => {
        return () => {
            if (debounceTimeoutRef.current) {
                clearTimeout(debounceTimeoutRef.current)
            }

            if (previewAbortControllerRef.current) {
                previewAbortControllerRef.current.abort()
            }
        }
    }, [])

    const handleStartTrip = async (fare: RouteFare) => {
        const payload = {
            rideFareId: fare.id,
            userId: userID,
        } as HTTPTripStartRequestPayload

        if (!fare.id) {
            alert("No Fare ID in the payload")
            return
        }

        const response = await fetch(`${API_URL}${BackendEndpoints.START_TRIP}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload),
        })
        const body = await response.json() as {
            data?: HTTPTripStartResponse
            error?: { message?: string }
        }
        if (!response.ok || !body.data) {
            throw new Error(body.error?.message || "Unable to start this trip right now")
        }

        return body.data
    }

    const handleCancelTrip = () => {
        setTrip(null)
        setDestination(null)
        setPreviewError(null)
        resetTripStatus()
    }

    if (error) {
        return <div>Error: {error}</div>
    }

    return (
        <div className="relative flex flex-col md:flex-row h-screen">
            <div className={`${destination ? 'flex-[0.7]' : 'flex-1'}`}>
                <MapContainer
                    center={[TEST_RIDER_LOCATION.latitude, TEST_RIDER_LOCATION.longitude]}
                    zoom={13}
                    style={{ height: '100%', width: '100%' }}
                    ref={mapRef}
                >
                    <TileLayer
                        url="https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png"
                        attribution="&copy; <a href='https://www.openstreetmap.org/copyright'>OpenStreetMap</a> contributors &copy; <a href='https://carto.com/'>CARTO</a>"
                    />
                    <Marker position={[TEST_RIDER_LOCATION.latitude, TEST_RIDER_LOCATION.longitude]} icon={userMarker} />

                    {/* Render geohash grid cells */}
                    {drivers?.map((driver) => (
                        <Rectangle
                            key={`grid-${driver?.geohash}`}
                            bounds={getGeohashBounds(driver?.geohash) as L.LatLngBoundsExpression}
                            pathOptions={{
                                color: '#3388ff',
                                weight: 1,
                                fillOpacity: 0.1
                            }}
                        >
                            <Popup>Geohash: {driver?.geohash}</Popup>
                        </Rectangle>
                    ))}

                    {/* Render driver markers */}
                    {drivers?.map((driver) => (
                        <Marker
                            key={driver?.id}
                            position={[driver?.location?.latitude, driver?.location?.longitude]}
                            icon={driverMarker}
                        >
                            <Popup>
                                Driver ID: {driver?.id}
                                <br />
                                Geohash: {driver?.geohash}
                                <br />
                                Name: {driver?.name}
                                <br />
                                Car Plate: {driver?.carPlate}
                                <br />
                                <DriverAvatar
                                    src={driver?.profilePicture}
                                    alt={`${driver?.name}'s profile picture`}
                                    width={100}
                                    height={100}
                                />
                            </Popup>
                        </Marker>
                    ))}
                    {destination && (
                        <Marker position={destination} icon={userMarker}>
                            <Popup>Destination</Popup>
                        </Marker>
                    )}

                    {selectedCarPackage && (
                        <div className="mt-4 z-[9999] absolute bottom-0 right-0">
                            <Button className="w-full">
                                Request Ride with {selectedCarPackage.packageSlug}
                            </Button>
                        </div>
                    )}
                    {trip && (
                        <RoutingControl route={trip.route} />
                    )}
                    <MapClickHandler onClick={handleMapClick} />
                </MapContainer>
            </div>

            <div className="flex-[0.4]">
                <RiderTripOverview
                    trip={trip}
                    previewError={previewError}
                    assignedDriver={assignedDriver}
                    status={tripStatus}
                    paymentSession={paymentSession}
                    onPackageSelect={handleStartTrip}
                    onCancel={handleCancelTrip}
                />
            </div>
        </div>
    )
}
