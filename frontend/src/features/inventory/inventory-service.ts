import { mockDealerships, mockVehicles } from './mock-data';
import type {
  CreateActionInput,
  Dealership,
  InventoryAction,
  InventoryFilters,
  InventoryStock,
  InventoryStockList,
  StockHistoryEvent,
} from './types';

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080';
const USE_API = process.env.NEXT_PUBLIC_INVENTORY_SOURCE === 'api';
const ACTION_STORAGE_KEY = 'keyloop-inventory-actions';

type ApiResponse<T> = {
  data: T;
};

function inventoryAgeDays(stockedAt: string): number {
  const today = new Date();
  today.setUTCHours(0, 0, 0, 0);
  const stocked = new Date(stockedAt);
  stocked.setUTCHours(0, 0, 0, 0);
  return Math.floor((today.getTime() - stocked.getTime()) / 86_400_000);
}

function toInventoryStock(vehicle: (typeof mockVehicles)[number]): InventoryStock {
  const age = inventoryAgeDays(vehicle.stockedAt);
  return {
    id: vehicle.id,
    dealershipId: vehicle.dealershipId,
    vehicleId: vehicle.id,
    vin: vehicle.vin,
    make: vehicle.make,
    model: vehicle.model,
    modelYear: vehicle.modelYear,
    price: vehicle.price,
    status: 'IN_STOCK',
    stockedInAt: vehicle.stockedAt,
    inventoryAgeDays: age,
    isAging: age > 90,
    latestAction: vehicle.latestAction
      ? {
          id: vehicle.latestAction.id,
          actionType: vehicle.latestAction.actionType,
          note: vehicle.latestAction.note,
          createdAt: vehicle.latestAction.createdAt,
        }
      : undefined,
  };
}

function getStoredActions(): Record<string, InventoryAction> {
  if (typeof window === 'undefined') return {};
  const stored = window.localStorage.getItem(ACTION_STORAGE_KEY);
  if (!stored) return {};
  try {
    return JSON.parse(stored) as Record<string, InventoryAction>;
  } catch {
    return {};
  }
}

function queryString(filters: InventoryFilters): string {
  const params = new URLSearchParams();
  if (filters.search) params.set('search', filters.search);
  if (filters.make) params.set('make', filters.make);
  if (filters.model) params.set('model', filters.model);
  if (filters.minAgeDays !== undefined) params.set('minAgeDays', String(filters.minAgeDays));
  if (filters.maxAgeDays !== undefined) params.set('maxAgeDays', String(filters.maxAgeDays));
  if (filters.agingOnly) params.set('agingOnly', 'true');
  if (filters.sortBy) params.set('sortBy', filters.sortBy);
  if (filters.sortOrder) params.set('sortOrder', filters.sortOrder);
  if (filters.page) params.set('page', String(filters.page));
  if (filters.pageSize) params.set('pageSize', String(filters.pageSize));
  return params.toString();
}

async function apiRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_URL}${path}`, {
    ...init,
    headers: { 'Content-Type': 'application/json', ...init?.headers },
  });
  if (!response.ok) throw new Error(`Inventory API request failed (${response.status})`);
  const result = (await response.json()) as ApiResponse<T>;
  return result.data;
}

export async function listDealerships(): Promise<Dealership[]> {
  if (USE_API) return apiRequest<Dealership[]>('/api/v1/dealerships');
  return Promise.resolve(mockDealerships);
}

export async function listInventoryStocks(
  dealershipId: string,
  filters: InventoryFilters,
): Promise<InventoryStockList> {
  if (USE_API) {
    const query = queryString(filters);
    return apiRequest<InventoryStockList>(
      `/api/v1/dealerships/${dealershipId}/stocks${query ? `?${query}` : ''}`,
    );
  }

  const storedActions = getStoredActions();
  const make = filters.make?.toLocaleLowerCase();
  const model = filters.model?.toLocaleLowerCase();

  const items = mockVehicles
    .filter((vehicle) => vehicle.dealershipId === dealershipId)
    .map((vehicle) => {
      const stock = toInventoryStock(vehicle);
      stock.latestAction = storedActions[vehicle.id] ?? stock.latestAction;
      return stock;
    })
    .filter((vehicle) => !make || vehicle.make.toLocaleLowerCase().includes(make))
    .filter((vehicle) => !model || vehicle.model.toLocaleLowerCase().includes(model))
    .filter((vehicle) => filters.minAgeDays === undefined || vehicle.inventoryAgeDays >= filters.minAgeDays)
    .filter((vehicle) => filters.maxAgeDays === undefined || vehicle.inventoryAgeDays <= filters.maxAgeDays)
    .filter((vehicle) => !filters.agingOnly || vehicle.isAging)
    .sort((left, right) => right.inventoryAgeDays - left.inventoryAgeDays);

  return {
    items,
    total: items.length,
    page: filters.page ?? 1,
    pageSize: filters.pageSize ?? items.length,
  };
}

export async function createStockAction(
  dealershipId: string,
  stockId: string,
  input: CreateActionInput,
): Promise<InventoryAction> {
  if (USE_API) {
    return apiRequest<InventoryAction>(
      `/api/v1/dealerships/${dealershipId}/stocks/${stockId}/actions`,
      { method: 'POST', body: JSON.stringify(input) },
    );
  }

  const action: InventoryAction = {
    id: crypto.randomUUID(),
    ...input,
    createdAt: new Date().toISOString(),
  };
  window.localStorage.setItem(
    ACTION_STORAGE_KEY,
    JSON.stringify({ ...getStoredActions(), [stockId]: action }),
  );
  return action;
}

export async function listStockHistory(
  dealershipId: string,
  stockId: string,
): Promise<StockHistoryEvent[]> {
  if (USE_API) {
    return apiRequest<StockHistoryEvent[]>(
      `/api/v1/dealerships/${dealershipId}/stocks/${stockId}/history`,
    );
  }

  const storedAction = getStoredActions()[stockId];
  const history: StockHistoryEvent[] = [
    {
      id: `${stockId}-stock-in`,
      eventType: 'STOCK_IN',
      note: 'Initial vehicle intake',
      occurredAt: new Date().toISOString(),
    },
  ];
  if (storedAction) {
    history.unshift({
      id: storedAction.id,
      eventType: 'ACTION',
      actionType: storedAction.actionType,
      note: storedAction.note,
      occurredAt: storedAction.createdAt,
    });
  }
  return history;
}
