package model

import "time"

type StockMovement struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)"`
	StockID      string    `gorm:"type:varchar(36);not null;index:stock_movements_stock_occurred_at_idx,priority:1"`
	MovementType string    `gorm:"type:varchar(20);not null;check:stock_movements_type_check,movement_type IN ('STOCK_IN','STOCK_OUT')"`
	Note         string    `gorm:"type:text;not null;check:stock_movements_note_length_check,char_length(note) BETWEEN 1 AND 500"`
	OccurredAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:stock_movements_stock_occurred_at_idx,priority:2,sort:desc"`
}

func (StockMovement) TableName() string { return "stock_movements" }
