import type { Dealership, InventoryAction } from './types';

type MockVehicle = {
  id: string;
  dealershipId: string;
  vin: string;
  make: string;
  model: string;
  modelYear: number;
  price: number;
  stockedAt: string;
  latestAction?: InventoryAction;
};

function dateDaysAgo(days: number): string {
  const date = new Date();
  date.setUTCHours(0, 0, 0, 0);
  date.setUTCDate(date.getUTCDate() - days);
  return date.toISOString();
}

export const mockDealerships: Dealership[] = [
  { id: 'dealer-hanoi', name: 'Keyloop Hanoi', location: 'Hanoi' },
  { id: 'dealer-saigon', name: 'Keyloop Saigon', location: 'Ho Chi Minh City' },
];

export const mockVehicles: MockVehicle[] = [
  {
    id: 'vehicle-001', dealershipId: 'dealer-hanoi', vin: 'KLHNX100000000001',
    make: 'Toyota', model: 'Camry', modelYear: 2024, price: 35800,
    stockedAt: dateDaysAgo(128),
    latestAction: {
      id: 'action-001', actionType: 'PRICE_REDUCTION_PLANNED',
      note: 'Review market pricing before the weekend campaign.', createdAt: dateDaysAgo(2),
    },
  },
  {
    id: 'vehicle-002', dealershipId: 'dealer-hanoi', vin: 'KLHNX100000000002',
    make: 'Honda', model: 'CR-V', modelYear: 2025, price: 32900,
    stockedAt: dateDaysAgo(96),
  },
  {
    id: 'vehicle-003', dealershipId: 'dealer-hanoi', vin: 'KLHNX100000000003',
    make: 'Mazda', model: 'CX-5', modelYear: 2025, price: 30100,
    stockedAt: dateDaysAgo(64),
  },
  {
    id: 'vehicle-004', dealershipId: 'dealer-hanoi', vin: 'KLHNX100000000004',
    make: 'Ford', model: 'Everest', modelYear: 2024, price: 44800,
    stockedAt: dateDaysAgo(181),
    latestAction: {
      id: 'action-002', actionType: 'TRANSFER_PROPOSED',
      note: 'Saigon location has stronger demand for this trim.', createdAt: dateDaysAgo(5),
    },
  },
  {
    id: 'vehicle-005', dealershipId: 'dealer-saigon', vin: 'KLSGN100000000001',
    make: 'Hyundai', model: 'Tucson', modelYear: 2025, price: 31500,
    stockedAt: dateDaysAgo(112),
  },
  {
    id: 'vehicle-006', dealershipId: 'dealer-saigon', vin: 'KLSGN100000000002',
    make: 'Toyota', model: 'Corolla Cross', modelYear: 2025, price: 28700,
    stockedAt: dateDaysAgo(43),
  },
  {
    id: 'vehicle-007', dealershipId: 'dealer-saigon', vin: 'KLSGN100000000003',
    make: 'Kia', model: 'Sportage', modelYear: 2024, price: 29600,
    stockedAt: dateDaysAgo(154),
  },
];
