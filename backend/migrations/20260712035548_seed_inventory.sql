-- Seed dealerships.
INSERT INTO dealerships (id, name, location) VALUES
  ('dealer-hanoi', 'Keyloop Hanoi', 'Hanoi'),
  ('dealer-saigon', 'Keyloop Saigon', 'Ho Chi Minh City');

-- Seed vehicle identity separately from dealership stock.
INSERT INTO vehicles (id, vin, make, model, model_year) VALUES
  ('vehicle-001', 'JTDBR32E720123456', 'Toyota', 'Corolla', 2022),
  ('vehicle-002', 'WVWZZZ1JZXW000001', 'Volkswagen', 'Golf', 2021),
  ('vehicle-003', '1HGCM82633A004352', 'Honda', 'Accord', 2023),
  ('vehicle-004', 'WF0AXXWPMAGR12345', 'Ford', 'Focus', 2020);

-- Relative dates keep aging examples meaningful whenever the stack starts.
INSERT INTO inventory_stocks (
  id, dealership_id, vehicle_id, status, price, stocked_in_at, stocked_out_at
) VALUES
  ('stock-001', 'dealer-hanoi', 'vehicle-001', 'IN_STOCK', 24500.00, CURRENT_TIMESTAMP - INTERVAL '120 days', NULL),
  ('stock-002', 'dealer-hanoi', 'vehicle-002', 'IN_STOCK', 21900.00, CURRENT_TIMESTAMP - INTERVAL '45 days', NULL),
  ('stock-003', 'dealer-hanoi', 'vehicle-003', 'OUT_OF_STOCK', 31500.00, CURRENT_TIMESTAMP - INTERVAL '150 days', CURRENT_TIMESTAMP - INTERVAL '5 days'),
  ('stock-004', 'dealer-saigon', 'vehicle-004', 'IN_STOCK', 18750.00, CURRENT_TIMESTAMP - INTERVAL '95 days', NULL);

-- Every seeded stock starts with an append-only stock-in event.
INSERT INTO stock_movements (id, stock_id, movement_type, note, occurred_at) VALUES
  ('movement-001', 'stock-001', 'STOCK_IN', 'Initial vehicle intake', CURRENT_TIMESTAMP - INTERVAL '120 days'),
  ('movement-002', 'stock-002', 'STOCK_IN', 'Initial vehicle intake', CURRENT_TIMESTAMP - INTERVAL '45 days'),
  ('movement-003', 'stock-003', 'STOCK_IN', 'Initial vehicle intake', CURRENT_TIMESTAMP - INTERVAL '150 days'),
  ('movement-004', 'stock-003', 'STOCK_OUT', 'Vehicle delivered to customer', CURRENT_TIMESTAMP - INTERVAL '5 days'),
  ('movement-005', 'stock-004', 'STOCK_IN', 'Initial vehicle intake', CURRENT_TIMESTAMP - INTERVAL '95 days');

INSERT INTO inventory_actions (id, stock_id, action_type, note, created_at) VALUES
  ('action-001', 'stock-001', 'PRICE_REDUCTION_PLANNED', 'Review a five percent reduction this week', CURRENT_TIMESTAMP - INTERVAL '2 days'),
  ('action-002', 'stock-004', 'TRANSFER_PROPOSED', 'Consider transfer to the Hanoi dealership', CURRENT_TIMESTAMP - INTERVAL '1 day');
