package model

import "time"

type InventoryStock struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)"`
	DealershipID string    `gorm:"type:varchar(36);not null;index:inventory_stocks_dealership_status_idx,priority:1;uniqueIndex:inventory_stocks_dealership_vehicle_unique,priority:1"`
	VehicleID    string    `gorm:"type:varchar(36);not null;uniqueIndex:inventory_stocks_dealership_vehicle_unique,priority:2"`
	Status       string    `gorm:"type:varchar(20);not null;index:inventory_stocks_dealership_status_idx,priority:2;check:inventory_stocks_status_check,status IN ('IN_STOCK','OUT_OF_STOCK')"`
	Price        float64   `gorm:"type:numeric(12,2);not null;check:inventory_stocks_price_check,price >= 0"`
	StockedInAt  time.Time `gorm:"not null;index"`
	StockedOutAt *time.Time
	CreatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

	Dealership Dealership        `gorm:"foreignKey:DealershipID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Vehicle    Vehicle           `gorm:"foreignKey:VehicleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Actions    []InventoryAction `gorm:"foreignKey:StockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Movements  []StockMovement   `gorm:"foreignKey:StockID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (InventoryStock) TableName() string { return "inventory_stocks" }
