import { Coordinate, Driver } from "../types";

export const OPEN_FREE_MAP_STYLE = "https://tiles.openfreemap.org/styles/liberty";

export interface GeoJsonFeatureCollection {
  type: "FeatureCollection";
  features: GeoJsonFeature[];
}

export interface GeoJsonFeature {
  type: "Feature";
  properties: Record<string, unknown>;
  geometry: {
    type: "LineString" | "Polygon";
    coordinates: number[][] | number[][][];
  };
}

export interface MapLibreGeoJsonSource {
  setData(data: GeoJsonFeatureCollection): void;
}

export interface MapLibreMap {
  addControl(control: unknown, position?: string): void;
  on(event: "load", handler: () => void): void;
  on(
    event: "click",
    handler: (event: { lngLat: { lat: number; lng: number } }) => void,
  ): void;
  getSource(id: string): MapLibreGeoJsonSource | undefined;
  addSource(id: string, source: { type: "geojson"; data: GeoJsonFeatureCollection }): void;
  addLayer(layer: {
    id: string;
    type: "fill" | "line";
    source: string;
    layout?: Record<string, string>;
    paint?: Record<string, string | number>;
  }): void;
  setCenter(center: [number, number]): void;
  remove(): void;
}

export interface MapLibrePopup {
  setDOMContent(node: Node): MapLibrePopup;
  setText(text: string): MapLibrePopup;
}

export interface MapLibreMarker {
  setLngLat(lngLat: [number, number]): MapLibreMarker;
  setPopup(popup: MapLibrePopup): MapLibreMarker;
  addTo(map: MapLibreMap): MapLibreMarker;
  remove(): void;
}

export interface MapLibreGlobal {
  Map: new (options: {
    container: HTMLElement;
    style: string;
    center: [number, number];
    zoom: number;
  }) => MapLibreMap;
  Marker: new (options: { element: HTMLElement; anchor: string }) => MapLibreMarker;
  Popup: new (options: { offset: number }) => MapLibrePopup;
  NavigationControl: new () => unknown;
}

declare global {
  interface Window {
    maplibregl?: MapLibreGlobal;
  }
}

export function getMapLibre(): MapLibreGlobal | null {
  return window.maplibregl ?? null;
}

export function coordinateToLngLat(coordinate: Coordinate): [number, number] {
  return [coordinate.longitude, coordinate.latitude];
}

export function tupleToLngLat([latitude, longitude]: [number, number]): [number, number] {
  return [longitude, latitude];
}

export function createCircleMarkerElement(label: string, color: string, size: number) {
  const element = document.createElement("div");
  element.style.width = `${size}px`;
  element.style.height = `${size}px`;
  element.style.borderRadius = "9999px";
  element.style.display = "flex";
  element.style.alignItems = "center";
  element.style.justifyContent = "center";
  element.style.background = color;
  element.style.color = "#ffffff";
  element.style.fontSize = `${Math.max(12, Math.floor(size / 2.2))}px`;
  element.style.fontWeight = "700";
  element.style.boxShadow = "0 10px 24px rgba(15, 23, 42, 0.18)";
  element.style.border = "2px solid rgba(255, 255, 255, 0.95)";
  element.textContent = label;
  return element;
}

export function createImageMarkerElement(iconUrl: string, width: number, height: number) {
  const element = document.createElement("div");
  element.style.width = `${width}px`;
  element.style.height = `${height}px`;
  element.style.backgroundImage = `url(${iconUrl})`;
  element.style.backgroundSize = "contain";
  element.style.backgroundRepeat = "no-repeat";
  element.style.backgroundPosition = "center";
  return element;
}

export function createDriverPopupContent(driver: Driver) {
  const container = document.createElement("div");
  container.style.minWidth = "180px";
  container.style.display = "grid";
  container.style.gap = "6px";
  container.style.fontSize = "13px";

  const lines = [
    `Driver ID: ${driver.id}`,
    `Geohash: ${driver.geohash}`,
    `Name: ${driver.name}`,
    `Car Plate: ${driver.carPlate}`,
  ];

  for (const line of lines) {
    const text = document.createElement("p");
    text.style.margin = "0";
    text.textContent = line;
    container.appendChild(text);
  }

  if (driver.profilePicture) {
    const image = document.createElement("img");
    image.src = driver.profilePicture;
    image.alt = `${driver.name}'s profile picture`;
    image.width = 100;
    image.height = 100;
    image.style.width = "100px";
    image.style.height = "100px";
    image.style.objectFit = "cover";
    image.style.borderRadius = "12px";
    image.style.marginTop = "4px";
    container.appendChild(image);
  }

  return container;
}

export function createTextPopupContent(lines: string[]) {
  const container = document.createElement("div");
  container.style.display = "grid";
  container.style.gap = "4px";
  container.style.fontSize = "13px";

  for (const line of lines) {
    const text = document.createElement("p");
    text.style.margin = "0";
    text.textContent = line;
    container.appendChild(text);
  }

  return container;
}
