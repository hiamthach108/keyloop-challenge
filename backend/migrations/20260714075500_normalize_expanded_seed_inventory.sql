-- Keep generated VINs clean and make expanded demo stock visible in the default dashboard.
UPDATE vehicles
SET vin = format('KL2026%s', lpad(g::text, 11, '0'))
FROM generate_series(5, 50) AS g
WHERE vehicles.id = format('vehicle-%s', lpad(g::text, 3, '0'));

UPDATE inventory_stocks
SET status = 'IN_STOCK', stocked_out_at = NULL
FROM generate_series(10, 50, 10) AS g
WHERE inventory_stocks.id = format('stock-%s', lpad(g::text, 3, '0'));

INSERT INTO stock_movements (id, stock_id, movement_type, note, occurred_at)
SELECT
  format('movement-%s-restock', lpad(g::text, 3, '0')),
  format('stock-%s', lpad(g::text, 3, '0')),
  'STOCK_IN',
  'Vehicle returned to active stock after status review',
  CURRENT_TIMESTAMP
FROM generate_series(10, 50, 10) AS g;
