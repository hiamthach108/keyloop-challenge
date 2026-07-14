-- Create "dealerships" table
CREATE TABLE "dealerships" (
  "id" character varying(36) NOT NULL,
  "name" character varying(255) NOT NULL,
  "location" character varying(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id")
);
-- Create "vehicles" table
CREATE TABLE "vehicles" (
  "id" character varying(36) NOT NULL,
  "vin" character varying(17) NOT NULL,
  "make" character varying(100) NOT NULL,
  "model" character varying(100) NOT NULL,
  "model_year" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "vehicles_model_year_check" CHECK (model_year >= 1900)
);
-- Create index "idx_vehicles_make" to table: "vehicles"
CREATE INDEX "idx_vehicles_make" ON "vehicles" ("make");
-- Create index "idx_vehicles_model" to table: "vehicles"
CREATE INDEX "idx_vehicles_model" ON "vehicles" ("model");
-- Create index "idx_vehicles_vin" to table: "vehicles"
CREATE UNIQUE INDEX "idx_vehicles_vin" ON "vehicles" ("vin");
-- Create "inventory_stocks" table
CREATE TABLE "inventory_stocks" (
  "id" character varying(36) NOT NULL,
  "dealership_id" character varying(36) NOT NULL,
  "vehicle_id" character varying(36) NOT NULL,
  "status" character varying(20) NOT NULL,
  "price" numeric(12,2) NOT NULL,
  "stocked_in_at" timestamptz NOT NULL,
  "stocked_out_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_dealerships_stocks" FOREIGN KEY ("dealership_id") REFERENCES "dealerships" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_inventory_stocks_vehicle" FOREIGN KEY ("vehicle_id") REFERENCES "vehicles" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "inventory_stocks_price_check" CHECK (price >= (0)::numeric),
  CONSTRAINT "inventory_stocks_status_check" CHECK ((status)::text = ANY ((ARRAY['IN_STOCK'::character varying, 'OUT_OF_STOCK'::character varying])::text[]))
);
-- Create index "idx_inventory_stocks_stocked_in_at" to table: "inventory_stocks"
CREATE INDEX "idx_inventory_stocks_stocked_in_at" ON "inventory_stocks" ("stocked_in_at");
-- Create index "inventory_stocks_dealership_status_idx" to table: "inventory_stocks"
CREATE INDEX "inventory_stocks_dealership_status_idx" ON "inventory_stocks" ("dealership_id", "status");
-- Create index "inventory_stocks_dealership_vehicle_unique" to table: "inventory_stocks"
CREATE UNIQUE INDEX "inventory_stocks_dealership_vehicle_unique" ON "inventory_stocks" ("dealership_id", "vehicle_id");
-- Create "inventory_actions" table
CREATE TABLE "inventory_actions" (
  "id" character varying(36) NOT NULL,
  "stock_id" character varying(36) NOT NULL,
  "action_type" character varying(50) NOT NULL,
  "note" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_inventory_stocks_actions" FOREIGN KEY ("stock_id") REFERENCES "inventory_stocks" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "inventory_actions_note_length_check" CHECK ((char_length(note) >= 1) AND (char_length(note) <= 500)),
  CONSTRAINT "inventory_actions_type_check" CHECK ((action_type)::text = ANY ((ARRAY['PRICE_REDUCTION_PLANNED'::character varying, 'TRANSFER_PROPOSED'::character varying, 'MARKETING_CAMPAIGN'::character varying, 'AWAITING_REVIEW'::character varying, 'OTHER'::character varying])::text[]))
);
-- Create index "inventory_actions_stock_created_at_idx" to table: "inventory_actions"
CREATE INDEX "inventory_actions_stock_created_at_idx" ON "inventory_actions" ("stock_id", "created_at" DESC);
-- Create "stock_movements" table
CREATE TABLE "stock_movements" (
  "id" character varying(36) NOT NULL,
  "stock_id" character varying(36) NOT NULL,
  "movement_type" character varying(20) NOT NULL,
  "note" text NOT NULL,
  "occurred_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_inventory_stocks_movements" FOREIGN KEY ("stock_id") REFERENCES "inventory_stocks" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "stock_movements_note_length_check" CHECK ((char_length(note) >= 1) AND (char_length(note) <= 500)),
  CONSTRAINT "stock_movements_type_check" CHECK ((movement_type)::text = ANY ((ARRAY['STOCK_IN'::character varying, 'STOCK_OUT'::character varying])::text[]))
);
-- Create index "stock_movements_stock_occurred_at_idx" to table: "stock_movements"
CREATE INDEX "stock_movements_stock_occurred_at_idx" ON "stock_movements" ("stock_id", "occurred_at" DESC);
