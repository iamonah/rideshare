"use client"

import { useDriverStreamConnection } from "../hooks/useDriverStreamConnection"
import { useMemo, useState } from "react";
import { useEffect, useRef } from "react";
import { CarPackageSlug, Coordinate } from "../types";
import { DriverTripOverview } from "./DriverTripOverview";
import * as Geohash from 'ngeohash';
import { DriverCard } from "./DriverCard";
import { TripEvents } from "../contracts";
import {
  coordinateToLngLat,
  createCircleMarkerElement,
  createImageMarkerElement,
  createTextPopupContent,
  GeoJsonFeatureCollection,
  getMapLibre,
  MapLibreMap,
  MapLibreMarker,
  OPEN_FREE_MAP_STYLE,
  tupleToLngLat,
} from "../lib/maplibre";

const START_LOCATION: Coordinate = {
  latitude: 37.7749,
  longitude: -122.4194,
}

export const DriverMap = ({ packageSlug }: { packageSlug: CarPackageSlug }) => {
  const mapContainerRef = useRef<HTMLDivElement | null>(null)
  const mapRef = useRef<MapLibreMap | null>(null)
  const markerRefs = useRef<MapLibreMarker[]>([])
  const userID = useMemo(() => crypto.randomUUID(), [])
  const [mapReady, setMapReady] = useState(false)
  const [riderLocation, setRiderLocation] = useState<Coordinate>(START_LOCATION)

  const driverGeohash = useMemo(() =>
    Geohash.encode(riderLocation?.latitude, riderLocation?.longitude, 7)
    , [riderLocation?.latitude, riderLocation?.longitude]);

  const {
    error,
    driver,
    tripStatus,
    requestedTrip,
    sendMessage,
    setTripStatus,
    resetTripStatus,
  } = useDriverStreamConnection({
    location: riderLocation,
    geohash: driverGeohash,
    userID,
    packageSlug,
  })

  const handleMapClick = (coordinate: Coordinate) => {
    setRiderLocation({
      latitude: coordinate.latitude,
      longitude: coordinate.longitude
    })
  }

  const handleAcceptTrip = () => {
    if (!requestedTrip || !requestedTrip.id || !driver) {
      alert("No trip ID found or driver is not set")
      return
    }

    sendMessage({
      type: TripEvents.DriverTripAccept,
      data: {
        tripID: requestedTrip.id,
        riderID: requestedTrip.userID,
        driver: driver,
      }
    })

    setTripStatus(TripEvents.DriverTripAccept)

  }

  const handleDeclineTrip = () => {
    if (!requestedTrip || !requestedTrip.id || !driver) {
      alert("No trip ID found or driver is not set")
      return
    }

    sendMessage({
      type: TripEvents.DriverTripDecline,
      data: {
        tripID: requestedTrip.id,
        riderID: requestedTrip.userID,
        driver: driver,
      }
    })

    setTripStatus(TripEvents.DriverTripDecline)
    resetTripStatus()
  }

  const parsedRoute = useMemo(() =>
    requestedTrip?.route?.geometry[0]?.coordinates
      .map((coord) => [coord?.latitude, coord?.longitude] as [number, number])
    , [requestedTrip])

  // destination is the last coordinate in the route
  const destination = useMemo(() =>
    requestedTrip?.route?.geometry[0]?.coordinates[requestedTrip?.route?.geometry[0]?.coordinates?.length - 1]
    , [requestedTrip])
  // start location is the first coordinate in the route
  const startLocation = useMemo(() =>
    requestedTrip?.route?.geometry[0]?.coordinates[0]
    , [requestedTrip])

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
        center: [START_LOCATION.longitude, START_LOCATION.latitude],
        zoom: 13,
      })

      map.addControl(new maplibregl.NavigationControl(), "top-right")
      map.on("load", () => setMapReady(true))
      map.on("click", (event: { lngLat: { lat: number; lng: number } }) => {
        handleMapClick({
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
    // The live location recentering happens in a separate effect.
  }, [])

  useEffect(() => {
    if (!mapReady || !mapRef.current) {
      return
    }

    mapRef.current.setCenter([riderLocation.longitude, riderLocation.latitude])
  }, [mapReady, riderLocation.latitude, riderLocation.longitude])

  useEffect(() => {
    if (!mapReady || !mapRef.current) {
      return
    }

    const map = mapRef.current
    const sourceId = "driver-route-source"
    const layerId = "driver-route-layer"
    const featureCollection: GeoJsonFeatureCollection = {
      type: "FeatureCollection",
      features: parsedRoute ? [{
        type: "Feature",
        geometry: {
          type: "LineString",
          coordinates: parsedRoute.map(tupleToLngLat),
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
  }, [mapReady, parsedRoute])

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
        element: createImageMarkerElement("https://www.svgrepo.com/show/25407/car.svg", 30, 30),
        anchor: "center",
      })
        .setLngLat([riderLocation.longitude, riderLocation.latitude])
        .setPopup(new maplibregl.Popup({ offset: 20 }).setDOMContent(createTextPopupContent([
          `Driver ID: ${userID}`,
          `Geohash: ${driverGeohash}`,
        ])))
        .addTo(map)
    )

    if (startLocation) {
      markerRefs.current.push(
        new maplibregl.Marker({
          element: createCircleMarkerElement("S", "#0f766e", 32),
          anchor: "center",
        })
          .setLngLat(coordinateToLngLat(startLocation))
          .setPopup(new maplibregl.Popup({ offset: 20 }).setText("Start Location"))
          .addTo(map)
      )
    }

    if (destination) {
      markerRefs.current.push(
        new maplibregl.Marker({
          element: createImageMarkerElement("/markers/location-pin.svg", 40, 40),
          anchor: "bottom",
        })
          .setLngLat(coordinateToLngLat(destination))
          .setPopup(new maplibregl.Popup({ offset: 20 }).setText("Destination"))
          .addTo(map)
      )
    }
  }, [destination, driverGeohash, mapReady, riderLocation.latitude, riderLocation.longitude, startLocation, userID])


  if (error) {
    return <div>Error: {error}</div>
  }

  return (
    <div className="relative flex flex-col md:flex-row h-screen">
      <div className="flex-1">
        <div ref={mapContainerRef} className="h-full w-full" />
        {!mapReady && (
          <div className="absolute inset-0 flex items-center justify-center bg-white/70 text-sm font-medium text-slate-700">
            Loading map...
          </div>
        )}
      </div>

      <div className="flex flex-col md:w-[400px] bg-white border-t md:border-t-0 md:border-l">
        <div className="p-4 border-b">
          <DriverCard driver={driver} packageSlug={packageSlug} />
        </div>
        <div className="flex-1 overflow-y-auto">
          <DriverTripOverview
            trip={requestedTrip}
            status={tripStatus}
            onAcceptTrip={handleAcceptTrip}
            onDeclineTrip={handleDeclineTrip}
          />
        </div>
      </div>
    </div>
  )
}
