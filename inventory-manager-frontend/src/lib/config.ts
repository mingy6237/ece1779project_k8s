// API base URL configuration
// - For minikube + Ingress: use empty string for same-origin requests via Ingress
// - For local development: defaults to http://localhost:8080
const apiBaseUrlEnv = process.env.NEXT_PUBLIC_API_BASE_URL;

// Use explicit env value if set (including empty string), otherwise use defaults
let resolvedApiBaseUrl: string;
if (apiBaseUrlEnv !== undefined) {
  resolvedApiBaseUrl = apiBaseUrlEnv; // Can be '', 'http://xxx', etc.
} else if (process.env.NODE_ENV === 'development') {
  // Local development: backend runs on host machine
  resolvedApiBaseUrl = 'http://localhost:8080';
} else {
  // Production/minikube: use same-origin via Ingress
  resolvedApiBaseUrl = '';
}

// Empty string = relative path (same-origin via Ingress)
export const API_BASE_URL = resolvedApiBaseUrl;

function deriveWsUrl() {
  // Use explicit WebSocket URL if configured
  if (process.env.NEXT_PUBLIC_WS_URL) {
    return process.env.NEXT_PUBLIC_WS_URL;
  }

  // If API_BASE_URL is empty, use same-origin WebSocket (recommended for Ingress)
  if (API_BASE_URL === '') {
    if (typeof window !== 'undefined') {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
      return `${protocol}//${window.location.host}/api/ws`;
    }
    // SSR/Node side: return relative path as placeholder
    return '/api/ws';
  }

  // If API_BASE_URL is a full URL, replace protocol to construct ws(s) URL
  if (API_BASE_URL.startsWith('https')) {
    return `${API_BASE_URL.replace('https', 'wss')}/api/ws`;
  }

  return `${API_BASE_URL.replace('http', 'ws')}/api/ws`;
}

export const WS_BASE_URL = deriveWsUrl();

// Server configuration
export interface ServerOption {
  id: string;
  name: string;
  url: string;
}

// Default server options:
// - First: same-origin Ingress (for minikube deployment)
// - Second: direct local backend connection (for development)
export const DEFAULT_SERVERS: ServerOption[] = [
  { id: 'ingress', name: 'Ingress (same origin)', url: '' },
  { id: 'local-8080', name: 'Local Dev 8080', url: 'http://localhost:8080' },
];