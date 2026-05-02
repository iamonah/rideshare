'use client';

import { useRiderStreamConnection } from '../hooks/useRiderStreamConnection';
import { getGeohashBounds } from '../utils/geohash';
import { useEffect, useMemo, useRef, useState } from 'react';
import { Button } from './ui/button';
import { RouteFare, RequestRideProps, TripPreview, HTTPTripStartResponse } from "../types";
import { API_URL } from '../constants';
import { RiderTripOverview } from './RiderTripOverview';
import { BackendEndpoints, HTTPTripPreviewRequestPayload, HTTPTripPreviewResponse, HTTPTripStartRequestPayload } from '../contracts';
import {
    coordinateToLngLat,
    createDriverPopupContent,
    createImageMarkerElement,
    GeoJsonFeatureCollection,
    getMapLibre,
    MapLibreMap,
    MapLibreMarker,
    OPEN_FREE_MAP_STYLE,
    tupleToLngLat,
} from '../lib/maplibre';

interface RiderMapProps {
    onRouteSelected?: (distance: number) => void;
}

const TEST_RIDER_LOCATION = {
    latitude: 37.7749,
    longitude: -122.4194,
}

export default function RiderMap({ onRouteSelected }: RiderMapProps) {
    const [trip, setTrip] = useState<TripPreview | null>(null)
    const [mapReady, setMapReady] = useState(false)
    const [previewError, setPreviewError] = useState<string | null>(null)
    const [selectedCarPackage] = useState<RouteFare | null>(null)
    const [destination, setDestination] = useState<[number, number] | null>(null)
    const mapContainerRef = useRef<HTMLDivElement | null>(null)
    const mapRef = useRef<MapLibreMap | null>(null)
    const markerRefs = useRef<MapLibreMarker[]>([])
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
    } = useRiderStreamConnection(userID);

    console.log(tripStatus)

    useEffect(() => {
        let cancelled = false
        let retryTimeout: number | null = null

        const initializeMap = () => {
            if (cancelled || mapRef.current || !mapContainerRef.current) {
                return
            }

            const maplibregl = getMapLibre()
            if (!maplibregl) {
                retryTimeout = window.setTimeout(initializeMap, 50)
                return
            }

            const map = new maplibregl.Map({
                container: mapContainerRef.current,
                style: OPEN_FREE_MAP_STYLE,
                center: [TEST_RIDER_LOCATION.longitude, TEST_RIDER_LOCATION.latitude],
                zoom: 13,
            })

            map.addControl(new maplibregl.NavigationControl(), "top-right")
            map.on("load", () => setMapReady(true))
            map.on("click", (event: { lngLat: { lat: number; lng: number } }) => {
                void handleMapClick({
                    latitude: event.lngLat.lat,
                    longitude: event.lngLat.lng,
                })
            })

            mapRef.current = map
        }

        initializeMap()

        return () => {
            cancelled = true
            if (retryTimeout) {
                clearTimeout(retryTimeout)
            }

            markerRefs.current.forEach((marker) => marker.remove())
            markerRefs.current = []

            if (mapRef.current) {
                mapRef.current.remove()
                mapRef.current = null
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        if (!mapReady || !mapRef.current) {
            return
        }

        mapRef.current.setCenter([TEST_RIDER_LOCATION.longitude, TEST_RIDER_LOCATION.latitude])
    }, [mapReady])

    const handleMapClick = async (coordinate: { latitude: number; longitude: number }) => {
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
            setDestination([coordinate.latitude, coordinate.longitude])
            setPreviewError(null)

            try {
                const data = await requestRidePreview({
                    pickup: [TEST_RIDER_LOCATION.latitude, TEST_RIDER_LOCATION.longitude],
                    destination: [coordinate.latitude, coordinate.longitude],
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

    useEffect(() => {
        if (!mapReady || !mapRef.current) {
            return
        }

        const map = mapRef.current
        const sourceId = "rider-geohash-source"
        const fillLayerId = "rider-geohash-fill"
        const lineLayerId = "rider-geohash-outline"
        const featureCollection: GeoJsonFeatureCollection = {
            type: "FeatureCollection",
            features: (drivers ?? [])
                .filter((driver) => Boolean(driver?.geohash))
                .map((driver) => {
                    const [[minLat, minLng], [maxLat, maxLng]] = getGeohashBounds(driver.geohash)

                    return {
                        type: "Feature",
                        properties: { geohash: driver.geohash },
                        geometry: {
                            type: "Polygon",
                            coordinates: [[
                                [minLng, minLat],
                                [maxLng, minLat],
                                [maxLng, maxLat],
                                [minLng, maxLat],
                                [minLng, minLat],
                            ]],
                        },
                    }
                }),
        }

        const source = map.getSource(sourceId)
        if (source) {
            source.setData(featureCollection)
            return
        }

        map.addSource(sourceId, {
            type: "geojson",
            data: featureCollection,
        })

        map.addLayer({
            id: fillLayerId,
            type: "fill",
            source: sourceId,
            paint: {
                "fill-color": "#3388ff",
                "fill-opacity": 0.08,
            },
        })

        map.addLayer({
            id: lineLayerId,
            type: "line",
            source: sourceId,
            paint: {
                "line-color": "#3388ff",
                "line-width": 1,
            },
        })
    }, [drivers, mapReady])

    useEffect(() => {
        if (!mapReady || !mapRef.current) {
            return
        }

        const map = mapRef.current
        const sourceId = "rider-route-source"
        const layerId = "rider-route-layer"
        const featureCollection: GeoJsonFeatureCollection = {
            type: "FeatureCollection",
            features: trip ? [{
                type: "Feature",
                geometry: {
                    type: "LineString",
                    coordinates: trip.route.map(tupleToLngLat),
                },
                properties: {},
            }] : [],
        }

        const source = map.getSource(sourceId)
        if (source) {
            source.setData(featureCollection)
            return
        }

        map.addSource(sourceId, {
            type: "geojson",
            data: featureCollection,
        })

        map.addLayer({
            id: layerId,
            type: "line",
            source: sourceId,
            layout: {
                "line-cap": "round",
                "line-join": "round",
            },
            paint: {
                "line-color": "#2563eb",
                "line-width": 5,
            },
        })
    }, [mapReady, trip])

    useEffect(() => {
        if (!mapReady || !mapRef.current) {
            return
        }

        const map = mapRef.current
        const maplibregl = getMapLibre()
        if (!maplibregl) {
            return
        }

        markerRefs.current.forEach((marker) => marker.remove())
        markerRefs.current = []

        markerRefs.current.push(
            new maplibregl.Marker({
                element: createImageMarkerElement("/markers/location-pin.svg", 40, 40),
                anchor: "bottom",
            })
                .setLngLat([TEST_RIDER_LOCATION.longitude, TEST_RIDER_LOCATION.latitude])
                .addTo(map)
        )

        for (const driver of drivers ?? []) {
            const popup = new maplibregl.Popup({ offset: 20 }).setDOMContent(createDriverPopupContent(driver))
            const marker = new maplibregl.Marker({
                element: createImageMarkerElement("https://www.svgrepo.com/show/25407/car.svg", 30, 30),
                anchor: "center",
            })
                .setLngLat(coordinateToLngLat(driver.location))
                .setPopup(popup)
                .addTo(map)

            markerRefs.current.push(marker)
        }

        if (destination) {
            markerRefs.current.push(
                new maplibregl.Marker({
                    element: createImageMarkerElement("/markers/location-pin.svg", 40, 40),
                    anchor: "bottom",
                })
                    .setLngLat(tupleToLngLat(destination))
                    .setPopup(new maplibregl.Popup({ offset: 20 }).setText("Destination"))
                    .addTo(map)
            )
        }
    }, [destination, drivers, mapReady])

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
                <div
                    ref={mapContainerRef}
                    className="h-full w-full"
                />
                {selectedCarPackage && (
                    <div className="mt-4 z-[9999] absolute bottom-0 right-0">
                        <Button className="w-full">
                            Request Ride with {selectedCarPackage.packageSlug}
                        </Button>
                    </div>
                )}
                {!mapReady && (
                    <div className="absolute inset-0 flex items-center justify-center bg-white/70 text-sm font-medium text-slate-700">
                        Loading map...
                    </div>
                )}
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
