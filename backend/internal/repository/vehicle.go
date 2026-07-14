package repository

import (
	"context"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"gorm.io/gorm"
)

type IVehicleRepository interface {
	FindByID(ctx context.Context, vehicleID string) (*model.Vehicle, error)
	FindByVIN(ctx context.Context, vin string) (*model.Vehicle, error)
}

type vehicleRepository struct{ db *gorm.DB }

func NewVehicleRepository(db *gorm.DB) IVehicleRepository {
	return &vehicleRepository{db: db}
}

func (r *vehicleRepository) FindByID(ctx context.Context, vehicleID string) (*model.Vehicle, error) {
	var row model.Vehicle
	if err := r.db.WithContext(ctx).Where("id = ?", vehicleID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *vehicleRepository) FindByVIN(ctx context.Context, vin string) (*model.Vehicle, error) {
	var row model.Vehicle
	if err := r.db.WithContext(ctx).Where("vin = ?", vin).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
