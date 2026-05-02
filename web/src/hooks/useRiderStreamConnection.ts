import { useEffect, useState } from 'react';
import { WEBSOCKET_URL } from "../constants";
import { Trip } from '../types';
import { Driver } from '../types';
import { PaymentEventSessionCreatedData, TripEvents, isValidWsMessage, BackendEndpoints, normalizeTrip, getWsMessageType } from '../contracts';

export function useRiderStreamConnection(userID: string) {
  const [drivers, setDrivers] = useState<Driver[]>([]);
  const [tripStatus, setTripStatus] = useState<TripEvents | null>(null);
  const [paymentSession, setPaymentSession] = useState<PaymentEventSessionCreatedData | null>(null);
  const [assignedDriver, setAssignedDriver] = useState<Trip["driver"] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!userID) return;

    const ws = new WebSocket(`${WEBSOCKET_URL}${BackendEndpoints.WS_RIDERS}?userID=${userID}`);

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data) as unknown;

      if (!message || !isValidWsMessage(message)) {
        setError(`Unknown message type "${getWsMessageType(message) ?? 'missing'}", allowed types are: ${Object.values(TripEvents).join(', ')}`);
        return;
      }

      setError(null);

      switch (message.type) {
        case TripEvents.Error:
          setError(message.error?.message ?? 'WebSocket server error');
          break;
        case TripEvents.DriverLocation:
          setDrivers(message.data);
          break;
        case TripEvents.PaymentSessionCreated:
          setPaymentSession(message.data);
          setTripStatus(message.type);
          break;
        case TripEvents.DriverAssigned:
          setAssignedDriver(normalizeTrip(message.data).driver);
          setTripStatus(message.type);
          break;
        case TripEvents.Created:
          setTripStatus(message.type);
          break;
        case TripEvents.NoDriversFound:
          setTripStatus(message.type);
          break;
      }
    };

    ws.onclose = () => {
      console.log('WebSocket closed');
    };

    ws.onerror = (event) => {
      setError('WebSocket error occurred');
      console.error('WebSocket error:', event);
    };

    return () => {
      console.log('Closing WebSocket');
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [userID]);

  const resetTripStatus = () => {
    setTripStatus(null);
    setPaymentSession(null);
  }

  return { drivers, assignedDriver, error, tripStatus, paymentSession, resetTripStatus };
}
