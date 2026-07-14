-- Expand the deterministic demo inventory to 50 stocks total.
INSERT INTO vehicles (id, vin, make, model, model_year)
SELECT
  format('vehicle-%s', lpad(g::text, 3, '0')),
  format('KL2026%011s', g::text),
  (ARRAY['Toyota', 'Honda', 'Ford', 'Hyundai', 'Kia', 'Mazda', 'Volkswagen', 'Nissan', 'BMW', 'Mercedes-Benz'])[((g - 5) % 10) + 1],
  (ARRAY['Corolla Cross', 'Civic', 'Ranger', 'Tucson', 'Sportage', 'CX-5', 'Tiguan', 'X-Trail', 'X3', 'C-Class'])[((g - 5) % 10) + 1],
  2020 + (g % 6)
FROM generate_series(5, 50) AS g;

INSERT INTO inventory_stocks (
  id, dealership_id, vehicle_id, status, price, stocked_in_at, stocked_out_at
)
SELECT
  format('stock-%s', lpad(g::text, 3, '0')),
  CASE WHEN g % 2 = 0 THEN 'dealer-hanoi' ELSE 'dealer-saigon' END,
  format('vehicle-%s', lpad(g::text, 3, '0')),
  CASE WHEN g % 10 = 0 THEN 'OUT_OF_STOCK' ELSE 'IN_STOCK' END,
  (18500 + (g * 875))::numeric(12, 2),
  CURRENT_TIMESTAMP - (((20 + (g * 7 % 180))::text || ' days')::interval),
  CASE
    WHEN g % 10 = 0 THEN CURRENT_TIMESTAMP - (((1 + (g % 12))::text || ' days')::interval)
    ELSE NULL
  END
FROM generate_series(5, 50) AS g;

INSERT INTO stock_movements (id, stock_id, movement_type, note, occurred_at)
SELECT
  format('movement-%s-in', lpad(g::text, 3, '0')),
  format('stock-%s', lpad(g::text, 3, '0')),
  'STOCK_IN',
  'Initial vehicle intake',
  CURRENT_TIMESTAMP - (((20 + (g * 7 % 180))::text || ' days')::interval)
FROM generate_series(5, 50) AS g;

INSERT INTO stock_movements (id, stock_id, movement_type, note, occurred_at)
SELECT
  format('movement-%s-out', lpad(g::text, 3, '0')),
  format('stock-%s', lpad(g::text, 3, '0')),
  'STOCK_OUT',
  'Vehicle sold and removed from active stock',
  CURRENT_TIMESTAMP - (((1 + (g % 12))::text || ' days')::interval)
FROM generate_series(5, 50) AS g
WHERE g % 10 = 0;

INSERT INTO inventory_actions (id, stock_id, action_type, note, created_at)
SELECT
  format('action-%s', lpad(g::text, 3, '0')),
  format('stock-%s', lpad(g::text, 3, '0')),
  CASE g % 5
    WHEN 0 THEN 'PRICE_REDUCTION_PLANNED'
    WHEN 1 THEN 'TRANSFER_PROPOSED'
    WHEN 2 THEN 'MARKETING_CAMPAIGN'
    WHEN 3 THEN 'AWAITING_REVIEW'
    ELSE 'OTHER'
  END,
  CASE g % 5
    WHEN 0 THEN 'Review price position against local market this week'
    WHEN 1 THEN 'Consider moving this vehicle to the other dealership'
    WHEN 2 THEN 'Add to next weekend campaign inventory list'
    WHEN 3 THEN 'Awaiting sales manager review before action'
    ELSE 'Custom follow-up required for this stock item'
  END,
  CURRENT_TIMESTAMP - (((1 + (g % 14))::text || ' days')::interval)
FROM generate_series(5, 50) AS g
WHERE g % 10 <> 0
  AND (20 + (g * 7 % 180)) > 90
  AND g % 3 <> 0;
