package repository

import (
	"context"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"gorm.io/gorm"
)

type IStockMovementRepository interface {
	ListByStockID(ctx context.Context, stockID string) ([]model.StockMovement, error)
}

type stockMovementRepository struct{ db *gorm.DB }

func NewStockMovementRepository(db *gorm.DB) IStockMovementRepository {
	return &stockMovementRepository{db: db}
}

func (r *stockMovementRepository) ListByStockID(ctx context.Context, stockID string) ([]model.StockMovement, error) {
	var rows []model.StockMovement
	err := r.db.WithContext(ctx).Where("stock_id = ?", stockID).Order("occurred_at DESC").Find(&rows).Error
	return rows, err
}
