package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hiamthach108/dreon-sdk/errorx"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/aggregate"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/repository"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/shared/constant"
	"github.com/hiamthach108/keyloop-challenge/backend/pkg/validator"
	"gorm.io/gorm"
)

type IInventorySvc interface {
	ListDealerships(ctx context.Context) ([]aggregate.DealershipAggregate, error)
	ListStocks(ctx context.Context, dealershipID string, filters aggregate.InventoryStockFilters) (*aggregate.InventoryStockListAggregate, error)
	CreateAction(ctx context.Context, dealershipID, stockID string, req *aggregate.CreateInventoryActionReq) (*aggregate.InventoryActionAggregate, error)
	RecordMovement(ctx context.Context, dealershipID, stockID string, req *aggregate.CreateStockMovementReq) (*aggregate.InventoryStockAggregate, error)
	GetHistory(ctx context.Context, dealershipID, stockID string) ([]aggregate.StockHistoryEventAggregate, error)
}

type InventorySvc struct {
	dealershipRepo repository.IDealershipRepository
	stockRepo      repository.IInventoryStockRepository
	actionRepo     repository.IInventoryActionRepository
	movementRepo   repository.IStockMovementRepository
	now            func() time.Time
}

func NewInventorySvc(
	dealershipRepo repository.IDealershipRepository,
	stockRepo repository.IInventoryStockRepository,
	actionRepo repository.IInventoryActionRepository,
	movementRepo repository.IStockMovementRepository,
) IInventorySvc {
	return &InventorySvc{
		dealershipRepo: dealershipRepo,
		stockRepo:      stockRepo,
		actionRepo:     actionRepo,
		movementRepo:   movementRepo,
		now:            time.Now,
	}
}

func (s *InventorySvc) ListDealerships(ctx context.Context) ([]aggregate.DealershipAggregate, error) {
	rows, err := s.dealershipRepo.List(ctx)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("list dealerships: %w", err))
	}
	result := make([]aggregate.DealershipAggregate, 0, len(rows))
	for i := range rows {
		var item aggregate.DealershipAggregate
		item.FromModel(&rows[i])
		result = append(result, item)
	}
	return result, nil
}

func (s *InventorySvc) ListStocks(ctx context.Context, dealershipID string, filters aggregate.InventoryStockFilters) (*aggregate.InventoryStockListAggregate, error) {
	if strings.TrimSpace(dealershipID) == "" {
		return nil, errorx.New(errorx.ErrBadRequest, "dealership ID is required")
	}
	if err := normalizeStockFilters(&filters); err != nil {
		return nil, err
	}
	exists, err := s.dealershipRepo.Exists(ctx, dealershipID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("find dealership: %w", err))
	}
	if !exists {
		return nil, errorx.New(errorx.ErrNotFound, "dealership not found")
	}

	now := s.now()
	rows, total, err := s.stockRepo.List(ctx, dealershipID, filters, now)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("list inventory stocks: %w", err))
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.ID)
	}
	actions, err := s.actionRepo.LatestByStockIDs(ctx, ids)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("list latest inventory actions: %w", err))
	}

	items := make([]aggregate.InventoryStockAggregate, 0, len(rows))
	for i := range rows {
		var action *model.InventoryAction
		if latest, ok := actions[rows[i].ID]; ok {
			action = &latest
		}
		var item aggregate.InventoryStockAggregate
		item.FromModel(&rows[i], action, now)
		items = append(items, item)
	}
	return &aggregate.InventoryStockListAggregate{
		Items: items, Total: total, Page: filters.Page, PageSize: filters.PageSize,
	}, nil
}

func (s *InventorySvc) CreateAction(ctx context.Context, dealershipID, stockID string, req *aggregate.CreateInventoryActionReq) (*aggregate.InventoryActionAggregate, error) {
	if err := validateIDs(dealershipID, stockID); err != nil {
		return nil, err
	}
	if err := validator.ValidateStruct(req); err != nil {
		return nil, errorx.Wrap(errorx.ErrBadRequest, validator.FormatValidationError(err))
	}
	if !req.ActionType.IsValid() {
		return nil, errorx.New(errorx.ErrBadRequest, "unsupported actionType")
	}
	stock, err := s.findStock(ctx, dealershipID, stockID)
	if err != nil {
		return nil, err
	}
	if constant.StockStatus(stock.Status) != constant.StockStatusInStock {
		return nil, errorx.New(errorx.ErrUnprocessable, "actions can only be logged for in-stock vehicles")
	}
	var current aggregate.InventoryStockAggregate
	current.FromModel(stock, nil, s.now())
	if !current.IsAging {
		return nil, errorx.New(errorx.ErrUnprocessable, "actions can only be logged for inventory older than 90 days")
	}

	action := &model.InventoryAction{
		ID: uuid.NewString(), StockID: stockID, ActionType: string(req.ActionType),
		Note: strings.TrimSpace(req.Note), CreatedAt: s.now().UTC(),
	}
	if err := s.actionRepo.Create(ctx, action); err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("create inventory action: %w", err))
	}
	var result aggregate.InventoryActionAggregate
	result.FromModel(action)
	return &result, nil
}

func (s *InventorySvc) RecordMovement(ctx context.Context, dealershipID, stockID string, req *aggregate.CreateStockMovementReq) (*aggregate.InventoryStockAggregate, error) {
	if err := validateIDs(dealershipID, stockID); err != nil {
		return nil, err
	}
	if err := validator.ValidateStruct(req); err != nil {
		return nil, errorx.Wrap(errorx.ErrBadRequest, validator.FormatValidationError(err))
	}
	if !req.MovementType.IsValid() {
		return nil, errorx.New(errorx.ErrBadRequest, "unsupported movementType")
	}
	stock, err := s.findStock(ctx, dealershipID, stockID)
	if err != nil {
		return nil, err
	}

	currentStatus := constant.StockStatus(stock.Status)
	expectedStatus := constant.StockStatusOutOfStock
	newStatus := constant.StockStatusInStock
	if req.MovementType == constant.StockMovementTypeOut {
		expectedStatus = constant.StockStatusInStock
		newStatus = constant.StockStatusOutOfStock
	}
	if currentStatus != expectedStatus {
		return nil, errorx.New(errorx.ErrConflict, "stock movement is invalid for the current status")
	}
	now := s.now().UTC()
	movement := &model.StockMovement{
		ID: uuid.NewString(), StockID: stockID, MovementType: string(req.MovementType),
		Note: strings.TrimSpace(req.Note), OccurredAt: now,
	}
	updated, err := s.stockRepo.RecordMovement(ctx, dealershipID, stockID, expectedStatus, newStatus, movement)
	if errors.Is(err, repository.ErrStockStateChanged) {
		return nil, errorx.New(errorx.ErrConflict, "stock status changed; reload and retry")
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("record stock movement: %w", err))
	}
	updated.Vehicle = stock.Vehicle
	var result aggregate.InventoryStockAggregate
	result.FromModel(updated, nil, now)
	return &result, nil
}

func (s *InventorySvc) GetHistory(ctx context.Context, dealershipID, stockID string) ([]aggregate.StockHistoryEventAggregate, error) {
	if err := validateIDs(dealershipID, stockID); err != nil {
		return nil, err
	}
	if _, err := s.findStock(ctx, dealershipID, stockID); err != nil {
		return nil, err
	}
	actions, err := s.actionRepo.ListByStockID(ctx, stockID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("list inventory actions: %w", err))
	}
	movements, err := s.movementRepo.ListByStockID(ctx, stockID)
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("list stock movements: %w", err))
	}
	return aggregate.StockHistoryFromModels(actions, movements), nil
}

func (s *InventorySvc) findStock(ctx context.Context, dealershipID, stockID string) (*model.InventoryStock, error) {
	stock, err := s.stockRepo.FindByID(ctx, dealershipID, stockID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorx.New(errorx.ErrNotFound, "inventory stock not found")
	}
	if err != nil {
		return nil, errorx.Wrap(errorx.ErrInternal, fmt.Errorf("find inventory stock: %w", err))
	}
	return stock, nil
}

func validateIDs(dealershipID, stockID string) error {
	if strings.TrimSpace(dealershipID) == "" || strings.TrimSpace(stockID) == "" {
		return errorx.New(errorx.ErrBadRequest, "dealership ID and stock ID are required")
	}
	return nil
}

func normalizeStockFilters(filters *aggregate.InventoryStockFilters) error {
	filters.Search = strings.TrimSpace(filters.Search)
	filters.Make = strings.TrimSpace(filters.Make)
	filters.Model = strings.TrimSpace(filters.Model)
	if filters.Status == "" {
		filters.Status = constant.StockStatusInStock
	}
	if !filters.Status.IsValid() {
		return errorx.New(errorx.ErrBadRequest, "unsupported status")
	}
	if filters.AgingOnly && filters.Status != constant.StockStatusInStock {
		return errorx.New(errorx.ErrBadRequest, "agingOnly requires IN_STOCK status")
	}
	if filters.MinAgeDays != nil && *filters.MinAgeDays < 0 {
		return errorx.New(errorx.ErrBadRequest, "minAgeDays must be non-negative")
	}
	if filters.MaxAgeDays != nil && *filters.MaxAgeDays < 0 {
		return errorx.New(errorx.ErrBadRequest, "maxAgeDays must be non-negative")
	}
	if filters.MinAgeDays != nil && filters.MaxAgeDays != nil && *filters.MinAgeDays > *filters.MaxAgeDays {
		return errorx.New(errorx.ErrBadRequest, "minAgeDays cannot exceed maxAgeDays")
	}
	if filters.SortBy == "" {
		filters.SortBy = aggregate.StockSortByStockedInAt
	}
	if !filters.SortBy.IsValid() {
		return errorx.New(errorx.ErrBadRequest, "unsupported sortBy")
	}
	if filters.SortOrder == "" {
		filters.SortOrder = aggregate.SortOrderDesc
	}
	if !filters.SortOrder.IsValid() {
		return errorx.New(errorx.ErrBadRequest, "unsupported sortOrder")
	}
	if filters.Page == 0 {
		filters.Page = 1
	}
	if filters.Page < 1 {
		return errorx.New(errorx.ErrBadRequest, "page must be positive")
	}
	if filters.PageSize == 0 {
		filters.PageSize = 20
	}
	if filters.PageSize < 1 || filters.PageSize > 100 {
		return errorx.New(errorx.ErrBadRequest, "pageSize must be between 1 and 100")
	}
	return nil
}
