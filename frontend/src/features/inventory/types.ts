export type Dealership = {
  id: string;
  name: string;
  location: string;
};

export type InventoryActionType =
  | 'PRICE_REDUCTION_PLANNED'
  | 'TRANSFER_PROPOSED'
  | 'MARKETING_CAMPAIGN'
  | 'AWAITING_REVIEW'
  | 'OTHER';

export type InventoryAction = {
  id: string;
  actionType: InventoryActionType;
  note: string;
  createdAt: string;
};

export type StockHistoryEventType = 'ACTION' | 'STOCK_IN' | 'STOCK_OUT';

export type StockHistoryEvent = {
  id: string;
  eventType: StockHistoryEventType;
  actionType?: InventoryActionType;
  note: string;
  occurredAt: string;
};

export type StockStatus = 'IN_STOCK' | 'OUT_OF_STOCK';

export type InventoryStock = {
  id: string;
  dealershipId: string;
  vehicleId: string;
  vin: string;
  make: string;
  model: string;
  modelYear: number;
  price: number;
  status: StockStatus;
  stockedInAt: string;
  stockedOutAt?: string;
  inventoryAgeDays: number;
  isAging: boolean;
  latestAction?: InventoryAction;
};

export type InventoryFilters = {
  search?: string;
  make?: string;
  model?: string;
  minAgeDays?: number;
  maxAgeDays?: number;
  agingOnly?: boolean;
  sortBy?: StockSortBy;
  sortOrder?: SortOrder;
  page?: number;
  pageSize?: number;
};

export type StockSortBy =
  | 'STOCKED_IN_AT'
  | 'MAKE'
  | 'MODEL'
  | 'MODEL_YEAR'
  | 'PRICE'
  | 'STATUS';

export type SortOrder = 'ASC' | 'DESC';

export type InventoryStockList = {
  items: InventoryStock[];
  total: number;
  page: number;
  pageSize: number;
};

export type CreateActionInput = {
  actionType: InventoryActionType;
  note: string;
};
