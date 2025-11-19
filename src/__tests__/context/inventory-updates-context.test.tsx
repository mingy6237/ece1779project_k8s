import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { InventoryUpdatesProvider, useInventoryUpdates } from '@/context/inventory-updates-context';
import { useAuth } from '@/context/auth-context';
import { InventoryUpdateEvent } from '@/lib/types';

// Mock the auth context
vi.mock('@/context/auth-context', () => ({
  useAuth: vi.fn(),
}));

// Mock WebSocket
class MockWebSocket {
  static instances: MockWebSocket[] = [];
  url: string;
  onopen: (() => void) | null = null;
  onclose: (() => void) | null = null;
  onerror: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  readyState: number = WebSocket.CONNECTING;

  constructor(url: string) {
    this.url = url;
    MockWebSocket.instances.push(this);
  }

  close() {
    this.readyState = WebSocket.CLOSED;
    if (this.onclose) {
      this.onclose();
    }
  }

  send(_data: string) {
    // Mock send - intentionally unused
  }

  // Helper method to simulate connection
  simulateOpen() {
    this.readyState = WebSocket.OPEN;
    if (this.onopen) {
      this.onopen();
    }
  }

  // Helper method to simulate message
  simulateMessage(data: string) {
    if (this.onmessage) {
      this.onmessage({ data });
    }
  }

  // Helper method to simulate error
  simulateError() {
    if (this.onerror) {
      this.onerror();
    }
  }

  static clearInstances() {
    MockWebSocket.instances = [];
  }
}

describe('InventoryUpdatesContext', () => {
  let originalWebSocket: typeof WebSocket;

  beforeEach(() => {
    // Store original WebSocket
    originalWebSocket = global.WebSocket;
    // Replace with mock
    global.WebSocket = MockWebSocket as unknown as typeof WebSocket;
    MockWebSocket.clearInstances();
    vi.clearAllMocks();
  });

  afterEach(() => {
    // Restore original WebSocket
    global.WebSocket = originalWebSocket;
    MockWebSocket.clearInstances();
  });

  it('should throw error when useInventoryUpdates is called outside provider', () => {
    // Suppress console.error for this test
    const originalError = console.error;
    console.error = vi.fn();

    expect(() => {
      renderHook(() => useInventoryUpdates());
    }).toThrow('useInventoryUpdates must be used within InventoryUpdatesProvider');

    console.error = originalError;
  });

  it('should not connect when token is null', () => {
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: null });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    expect(result.current.connected).toBe(false);
    expect(result.current.lastEvent).toBe(null);
    expect(MockWebSocket.instances.length).toBe(0);
  });

  it('should establish WebSocket connection when token is provided', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    expect(MockWebSocket.instances.length).toBe(1);
    const ws = MockWebSocket.instances[0];
    expect(ws.url).toContain('token=test-token-123');

    // Initially not connected
    expect(result.current.connected).toBe(false);

    // Simulate connection
    act(() => {
      ws.simulateOpen();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(true);
    });
  });

  it('should parse and store inventory update events', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    const mockEvent: InventoryUpdateEvent = {
      inventory_id: 'inv-123',
      sku_id: 'sku-456',
      sku_name: 'Test Product',
      store_id: 'store-789',
      store_name: 'Test Store',
      operation_type: 'adjust',
      old_quantity: 10,
      new_quantity: 15,
      delta_quantity: 5,
      version: 2,
      updated_at: '2025-11-18T00:00:00Z',
      updated_by: 'user-123',
    };

    act(() => {
      ws.simulateMessage(JSON.stringify(mockEvent));
    });

    await waitFor(() => {
      expect(result.current.lastEvent).toEqual(mockEvent);
    });
  });

  it('should handle multiple events and keep only the last one', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    const event1: InventoryUpdateEvent = {
      id: 'event-1',
      inventory_id: 'inv-1',
      sku_id: 'sku-1',
      sku_name: 'Product 1',
      store_id: 'store-1',
      store_name: 'Store 1',
      operation_type: 'create',
      sender_instance_id: 'instance-1',
      user_id: 'user-1',
      user_name: 'User 1',
      new_quantity: 10,
      delta_quantity: 10,
      version: 1,
      created_at: '2025-11-18T00:00:00Z',
      updated_at: '2025-11-18T00:00:00Z',
    };

    const event2: InventoryUpdateEvent = {
      id: 'event-2',
      inventory_id: 'inv-2',
      sku_id: 'sku-2',
      sku_name: 'Product 2',
      store_id: 'store-2',
      store_name: 'Store 2',
      operation_type: 'update',
      sender_instance_id: 'instance-1',
      user_id: 'user-2',
      user_name: 'User 2',
      new_quantity: 8,
      delta_quantity: 3,
      version: 2,
      created_at: '2025-11-18T01:00:00Z',
      updated_at: '2025-11-18T01:00:00Z',
    };

    act(() => {
      ws.simulateMessage(JSON.stringify(event1));
    });

    await waitFor(() => {
      expect(result.current.lastEvent).toEqual(event1);
    });

    act(() => {
      ws.simulateMessage(JSON.stringify(event2));
    });

    await waitFor(() => {
      expect(result.current.lastEvent).toEqual(event2);
    });
  });

  it('should clear last event when clearLastEvent is called', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    const mockEvent: InventoryUpdateEvent = {
      id: 'event-delete',
      inventory_id: 'inv-123',
      sku_id: 'sku-456',
      sku_name: 'Test Product',
      store_id: 'store-789',
      store_name: 'Test Store',
      operation_type: 'delete',
      sender_instance_id: 'instance-1',
      user_id: 'user-123',
      user_name: 'Test User',
      new_quantity: 0,
      delta_quantity: -5,
      version: 3,
      created_at: '2025-11-18T00:00:00Z',
      updated_at: '2025-11-18T00:00:00Z',
    };

    act(() => {
      ws.simulateMessage(JSON.stringify(mockEvent));
    });

    await waitFor(() => {
      expect(result.current.lastEvent).toEqual(mockEvent);
    });

    act(() => {
      result.current.clearLastEvent();
    });

    await waitFor(() => {
      expect(result.current.lastEvent).toBe(null);
    });
  });

  it('should handle invalid JSON messages gracefully', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const consoleErrorSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    act(() => {
      ws.simulateMessage('invalid json {{{');
    });

    await waitFor(() => {
      expect(consoleErrorSpy).toHaveBeenCalledWith(
        'Failed to parse inventory update',
        expect.any(Error)
      );
    });

    expect(result.current.lastEvent).toBe(null);

    consoleErrorSpy.mockRestore();
  });

  it('should set connected to false on WebSocket error', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(true);
    });

    act(() => {
      ws.simulateError();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(false);
    });
  });

  it('should set connected to false on WebSocket close', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(true);
    });

    act(() => {
      ws.close();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(false);
    });
  });

  it('should close WebSocket and clear state when token becomes null', async () => {
    const mockToken = 'test-token-123';
    const mockUseAuth = vi.fn().mockReturnValue({ token: mockToken });
    (useAuth as ReturnType<typeof vi.fn>).mockImplementation(mockUseAuth);

    const { result, rerender } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws1 = MockWebSocket.instances[0];
    
    act(() => {
      ws1.simulateOpen();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(true);
    });

    const mockEvent: InventoryUpdateEvent = {
      id: 'event-token-test',
      inventory_id: 'inv-123',
      sku_id: 'sku-456',
      sku_name: 'Test Product',
      store_id: 'store-789',
      store_name: 'Test Store',
      operation_type: 'adjust',
      sender_instance_id: 'instance-1',
      user_id: 'user-123',
      user_name: 'Test User',
      new_quantity: 15,
      delta_quantity: 5,
      version: 2,
      created_at: '2025-11-18T00:00:00Z',
      updated_at: '2025-11-18T00:00:00Z',
    };

    act(() => {
      ws1.simulateMessage(JSON.stringify(mockEvent));
    });

    await waitFor(() => {
      expect(result.current.lastEvent).toEqual(mockEvent);
    });

    // Change token to null
    mockUseAuth.mockReturnValue({ token: null });
    rerender();

    await waitFor(() => {
      expect(result.current.connected).toBe(false);
      expect(result.current.lastEvent).toBe(null);
    });
  });

  it('should create new WebSocket connection when token changes', async () => {
    const mockToken1 = 'token-1';
    const mockUseAuth = vi.fn().mockReturnValue({ token: mockToken1 });
    (useAuth as ReturnType<typeof vi.fn>).mockImplementation(mockUseAuth);

    const { result, rerender } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    expect(MockWebSocket.instances.length).toBe(1);
    const ws1 = MockWebSocket.instances[0];
    expect(ws1.url).toContain('token=token-1');

    act(() => {
      ws1.simulateOpen();
    });

    await waitFor(() => {
      expect(result.current.connected).toBe(true);
    });

    // Change token
    const mockToken2 = 'token-2';
    mockUseAuth.mockReturnValue({ token: mockToken2 });
    rerender();

    await waitFor(() => {
      expect(MockWebSocket.instances.length).toBe(2);
    });

    const ws2 = MockWebSocket.instances[1];
    expect(ws2.url).toContain('token=token-2');
  });

  it('should handle different operation types correctly', async () => {
    const mockToken = 'test-token-123';
    (useAuth as ReturnType<typeof vi.fn>).mockReturnValue({ token: mockToken });

    const { result } = renderHook(() => useInventoryUpdates(), {
      wrapper: InventoryUpdatesProvider,
    });

    const ws = MockWebSocket.instances[0];
    
    act(() => {
      ws.simulateOpen();
    });

    const operationTypes: Array<'create' | 'update' | 'adjust' | 'delete'> = [
      'create',
      'update',
      'adjust',
      'delete',
    ];

    for (const opType of operationTypes) {
      const event: InventoryUpdateEvent = {
        id: `event-${opType}`,
        inventory_id: `inv-${opType}`,
        sku_id: 'sku-1',
        sku_name: 'Test Product',
        store_id: 'store-1',
        store_name: 'Test Store',
        operation_type: opType,
        sender_instance_id: 'instance-1',
        user_id: 'user-1',
        user_name: 'User 1',
        new_quantity: opType === 'delete' ? 0 : 15,
        delta_quantity: opType === 'delete' ? -10 : 5,
        version: 1,
        created_at: '2025-11-18T00:00:00Z',
        updated_at: '2025-11-18T00:00:00Z',
      };

      act(() => {
        ws.simulateMessage(JSON.stringify(event));
      });

      await waitFor(() => {
        expect(result.current.lastEvent?.operation_type).toBe(opType);
      });
    }
  });
});
