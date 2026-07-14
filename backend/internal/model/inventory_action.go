package model

import "time"

type InventoryAction struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)"`
	StockID    string    `gorm:"type:varchar(36);not null;index:inventory_actions_stock_created_at_idx,priority:1"`
	ActionType string    `gorm:"type:varchar(50);not null;check:inventory_actions_type_check,action_type IN ('PRICE_REDUCTION_PLANNED','TRANSFER_PROPOSED','MARKETING_CAMPAIGN','AWAITING_REVIEW','OTHER')"`
	Note       string    `gorm:"type:text;not null;check:inventory_actions_note_length_check,char_length(note) BETWEEN 1 AND 500"`
	CreatedAt  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;index:inventory_actions_stock_created_at_idx,priority:2,sort:desc"`
}

func (InventoryAction) TableName() string { return "inventory_actions" }
