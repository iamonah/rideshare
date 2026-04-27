import { useEffect, useRef, useState } from 'react';
import { WEBSOCKET_URL } from "../constants";
import { Trip, Driver, CarPackageSlug } from '../types';
import { ServerWsMessage, TripEvents, isValidWsMessage, isValidTripEvent, ClientWsMessage, BackendEndpoints, normalizeDriver, normalizeTrip } from '../contracts';

interface useDriverConnectionProps {
  location: {
    latitude: number;
    longitude: number;
  };
  geohash: string;
  userID: string;
  packageSlug: CarPackageSlug;
}

export const useDriverStreamConnection = ({
  location,
  geohash,
  userID,
  packageSlug
}: useDriverConnectionProps) => {
  const [requestedTrip, setRequestedTrip] = useState<Trip | null>(null)
  const [tripStatus, setTripStatus] = useState<TripEvents | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [driver, setDriver] = useState<Driver | null>(null);
  const socketRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!userID) return;

    const websocket = new WebSocket(`${WEBSOCKET_URL}${BackendEndpoints.WS_DRIVERS}?userID=${userID}&packageSlug=${packageSlug}`);
    socketRef.current = websocket;

    websocket.onopen = () => {
      setError(null);

      if (location) {
        // Send initial location
        websocket.send(JSON.stringify({
          type: TripEvents.DriverLocation,
          data: {
            location,
            geohash,
          }
        }));
      }
    };

    websocket.onmessage = (event) => {
      const message = JSON.parse(event.data) as ServerWsMessage;

      if (!message || !isValidWsMessage(message)) {
        setError(`Unknown message type "${message}", allowed types are: ${Object.values(TripEvents).join(', ')}`);
        return;
      }

      switch (message.type) {
        case TripEvents.DriverTripRequest:
          setRequestedTrip(normalizeTrip(message.data));
          break;
        case TripEvents.DriverRegister:
          setDriver(normalizeDriver(message.data));
          break;
      }


      if (isValidTripEvent(message.type)) {
        setTripStatus(message.type);
      } else {
        setError(`Unknown message type "${message.type}", allowed types are: ${Object.values(TripEvents).join(', ')}`);
      }
    };

    websocket.onclose = (event) => {
      socketRef.current = null;
      console.log(`WebSocket closed (code: ${event.code}, reason: ${event.reason || 'no reason provided'})`);
    };

    websocket.onerror = (event) => {
      setError('WebSocket error occurred');
      console.error('WebSocket error:', event);
    };

    return () => {
      console.log('Closing WebSocket');
      socketRef.current = null;
      if (websocket.readyState === WebSocket.OPEN) {
        websocket.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID]);

  const sendMessage = (message: ClientWsMessage) => {
    const websocket = socketRef.current;

    if (websocket?.readyState === WebSocket.OPEN) {
      websocket.send(JSON.stringify(message));
    } else {
      const state = websocket ? socketStateLabel(websocket.readyState) : 'closed';
      setError(`WebSocket is not connected (state: ${state})`);
    }
  };

  const resetTripStatus = () => {
    setTripStatus(null);
    setRequestedTrip(null);
  }

  return { error, tripStatus, driver, requestedTrip, resetTripStatus, sendMessage, setTripStatus };
}

function socketStateLabel(state: number) {
  switch (state) {
    case WebSocket.CONNECTING:
      return 'connecting';
    case WebSocket.OPEN:
      return 'open';
    case WebSocket.CLOSING:
      return 'closing';
    case WebSocket.CLOSED:
      return 'closed';
    default:
      return `unknown(${state})`;
  }
}
