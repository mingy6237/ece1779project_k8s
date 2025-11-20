'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import { InventoryUpdateEvent } from '@/lib/types';
import { useAuth } from '@/context/auth-context';
import { useServer } from '@/context/server-context';
import { WS_BASE_URL } from '@/lib/config';

interface InventoryUpdatesContextValue {
  connected: boolean;
  lastEvent: InventoryUpdateEvent | null;
  clearLastEvent: () => void;
}

const InventoryUpdatesContext = createContext<InventoryUpdatesContextValue | undefined>(undefined);

export function InventoryUpdatesProvider({ children }: { children: React.ReactNode }) {
  const { token } = useAuth();
  const { selectedServer } = useServer();
  const [connected, setConnected] = useState(false);
  const [lastEvent, setLastEvent] = useState<InventoryUpdateEvent | null>(null);

  useEffect(() => {
    if (!token) {
      setConnected(false);
      setLastEvent(null);
      return;
    }

    // Derive WebSocket URL from selected server
    // If selectedServer.url is empty, use WS_BASE_URL (which handles relative paths)
    let wsBaseUrl: string;
    if (selectedServer.url === '') {
      // Use the configured WS_BASE_URL which handles relative paths correctly
      wsBaseUrl = WS_BASE_URL;
    } else {
      // Build from selected server URL
      wsBaseUrl = selectedServer.url.startsWith('https')
        ? `${selectedServer.url.replace('https', 'wss')}/api/ws`
        : `${selectedServer.url.replace('http', 'ws')}/api/ws`;
    }

    let wsUrl: URL;
    try {
      // If wsBaseUrl is a relative path (starts with /), construct full URL
      if (wsBaseUrl.startsWith('/')) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        wsUrl = new URL(`${protocol}//${window.location.host}${wsBaseUrl}`);
      } else {
        wsUrl = new URL(wsBaseUrl);
      }
    } catch (error) {
      console.error('Invalid WS URL', error, { wsBaseUrl, selectedServerUrl: selectedServer.url });
      return;
    }
    wsUrl.searchParams.set('token', token);

    const socket = new WebSocket(wsUrl.toString());

    socket.onopen = () => {
      console.log('WebSocket connected');
      setConnected(true);
    };
    socket.onclose = (event) => {
      console.log('WebSocket closed', event.code, event.reason);
      setConnected(false);
    };
    socket.onerror = (error) => {
      console.error('WebSocket error', error);
      setConnected(false);
    };
    socket.onmessage = (event) => {
      try {
        console.log('WebSocket message received:', event.data);
        const payload = JSON.parse(event.data) as InventoryUpdateEvent;
        console.log('Parsed inventory update event:', payload);
        setLastEvent(payload);
      } catch (error) {
        console.error('Failed to parse inventory update', error, event.data);
      }
    };

    return () => {
      socket.close();
    };
  }, [token, selectedServer.url]);

  const clearLastEvent = () => setLastEvent(null);

  return (
    <InventoryUpdatesContext.Provider value={{ connected, lastEvent, clearLastEvent }}>
      {children}
    </InventoryUpdatesContext.Provider>
  );
}

export function useInventoryUpdates() {
  const context = useContext(InventoryUpdatesContext);
  if (!context) {
    throw new Error('useInventoryUpdates must be used within InventoryUpdatesProvider');
  }
  return context;
}

