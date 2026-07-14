package model

import "time"

type Dealership struct {
	ID        string           `gorm:"primaryKey;type:varchar(36)"`
	Name      string           `gorm:"type:varchar(255);not null"`
	Location  string           `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time        `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Stocks    []InventoryStock `gorm:"foreignKey:DealershipID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Dealership) TableName() string { return "dealerships" }
