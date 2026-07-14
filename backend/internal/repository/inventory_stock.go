package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/aggregate"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/shared/constant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrStockStateChanged = errors.New("inventory stock state changed")

type IInventoryStockRepository interface {
	List(ctx context.Context, dealershipID string, filters aggregate.InventoryStockFilters, now time.Time) ([]model.InventoryStock, int64, error)
	FindByID(ctx context.Context, dealershipID, stockID string) (*model.InventoryStock, error)
	RecordMovement(ctx context.Context, dealershipID, stockID string, expectedStatus, newStatus constant.StockStatus, movement *model.StockMovement) (*model.InventoryStock, error)
}

type inventoryStockRepository struct{ db *gorm.DB }

func NewInventoryStockRepository(db *gorm.DB) IInventoryStockRepository {
	return &inventoryStockRepository{db: db}
}

var stockSortColumns = map[aggregate.StockSortBy]string{
	aggregate.StockSortByStockedInAt: "inventory_stocks.stocked_in_at",
	aggregate.StockSortByMake:        `"Vehicle"."make"`,
	aggregate.StockSortByModel:       `"Vehicle"."model"`,
	aggregate.StockSortByModelYear:   `"Vehicle"."model_year"`,
	aggregate.StockSortByPrice:       "inventory_stocks.price",
	aggregate.StockSortByStatus:      "inventory_stocks.status",
}

func (r *inventoryStockRepository) List(ctx context.Context, dealershipID string, filters aggregate.InventoryStockFilters, now time.Time) ([]model.InventoryStock, int64, error) {
	query := r.db.WithContext(ctx).Model(&model.InventoryStock{}).
		Joins("Vehicle").
		Where("inventory_stocks.dealership_id = ?", dealershipID)
	if filters.Search != "" {
		pattern := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where(`LOWER("Vehicle"."make") LIKE ? OR LOWER("Vehicle"."model") LIKE ?`, pattern, pattern)
	}
	if filters.Make != "" {
		query = query.Where(`LOWER("Vehicle"."make") LIKE ?`, "%"+strings.ToLower(filters.Make)+"%")
	}
	if filters.Model != "" {
		query = query.Where(`LOWER("Vehicle"."model") LIKE ?`, "%"+strings.ToLower(filters.Model)+"%")
	}
	query = query.Where("inventory_stocks.status = ?", filters.Status)
	today := time.Date(now.UTC().Year(), now.UTC().Month(), now.UTC().Day(), 0, 0, 0, 0, time.UTC)
	if filters.MinAgeDays != nil {
		query = query.Where("inventory_stocks.stocked_in_at <= ?", today.AddDate(0, 0, -*filters.MinAgeDays))
	}
	if filters.MaxAgeDays != nil {
		query = query.Where("inventory_stocks.stocked_in_at >= ?", today.AddDate(0, 0, -*filters.MaxAgeDays))
	}
	if filters.AgingOnly {
		query = query.Where("inventory_stocks.stocked_in_at < ?", today.AddDate(0, 0, -90))
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	order := fmt.Sprintf("%s %s, inventory_stocks.id ASC", stockSortColumns[filters.SortBy], filters.SortOrder)
	var rows []model.InventoryStock
	err := query.Order(order).Limit(filters.PageSize).Offset((filters.Page - 1) * filters.PageSize).Find(&rows).Error
	return rows, total, err
}

func (r *inventoryStockRepository) FindByID(ctx context.Context, dealershipID, stockID string) (*model.InventoryStock, error) {
	var row model.InventoryStock
	err := r.db.WithContext(ctx).Preload("Vehicle").
		Where("inventory_stocks.id = ? AND inventory_stocks.dealership_id = ?", stockID, dealershipID).
		First(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *inventoryStockRepository) RecordMovement(ctx context.Context, dealershipID, stockID string, expectedStatus, newStatus constant.StockStatus, movement *model.StockMovement) (*model.InventoryStock, error) {
	var stock model.InventoryStock
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ? AND dealership_id = ?", stockID, dealershipID).First(&stock).Error; err != nil {
			return err
		}
		if constant.StockStatus(stock.Status) != expectedStatus {
			return ErrStockStateChanged
		}
		updates := map[string]any{"status": newStatus, "updated_at": movement.OccurredAt}
		if newStatus == constant.StockStatusInStock {
			updates["stocked_in_at"] = movement.OccurredAt
			updates["stocked_out_at"] = nil
		} else {
			updates["stocked_out_at"] = movement.OccurredAt
		}
		if err := tx.Model(&stock).Updates(updates).Error; err != nil {
			return err
		}
		return tx.Create(movement).Error
	})
	if err != nil {
		return nil, err
	}
	stock.Status = string(newStatus)
	return &stock, nil
}
