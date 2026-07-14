package repository

import (
	"context"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"gorm.io/gorm"
)

type IInventoryActionRepository interface {
	Create(ctx context.Context, action *model.InventoryAction) error
	LatestByStockIDs(ctx context.Context, stockIDs []string) (map[string]model.InventoryAction, error)
	ListByStockID(ctx context.Context, stockID string) ([]model.InventoryAction, error)
}

type inventoryActionRepository struct{ db *gorm.DB }

func NewInventoryActionRepository(db *gorm.DB) IInventoryActionRepository {
	return &inventoryActionRepository{db: db}
}

func (r *inventoryActionRepository) Create(ctx context.Context, action *model.InventoryAction) error {
	return r.db.WithContext(ctx).Create(action).Error
}

func (r *inventoryActionRepository) LatestByStockIDs(ctx context.Context, stockIDs []string) (map[string]model.InventoryAction, error) {
	result := make(map[string]model.InventoryAction)
	if len(stockIDs) == 0 {
		return result, nil
	}
	var rows []model.InventoryAction
	err := r.db.WithContext(ctx).Where("stock_id IN ?", stockIDs).
		Order("stock_id ASC, created_at DESC").Find(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if _, exists := result[row.StockID]; !exists {
			result[row.StockID] = row
		}
	}
	return result, nil
}

func (r *inventoryActionRepository) ListByStockID(ctx context.Context, stockID string) ([]model.InventoryAction, error) {
	var rows []model.InventoryAction
	err := r.db.WithContext(ctx).Where("stock_id = ?", stockID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}
