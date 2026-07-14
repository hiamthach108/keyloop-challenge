package model

import "time"

type Vehicle struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)"`
	VIN       string    `gorm:"type:varchar(17);not null;uniqueIndex"`
	Make      string    `gorm:"type:varchar(100);not null;index"`
	Model     string    `gorm:"type:varchar(100);not null;index"`
	ModelYear int       `gorm:"not null;check:vehicles_model_year_check,model_year >= 1900"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

func (Vehicle) TableName() string { return "vehicles" }
