package aggregate

import (
	"time"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/shared/constant"
)

type StockSortBy string

const (
	StockSortByStockedInAt StockSortBy = "STOCKED_IN_AT"
	StockSortByMake        StockSortBy = "MAKE"
	StockSortByModel       StockSortBy = "MODEL"
	StockSortByModelYear   StockSortBy = "MODEL_YEAR"
	StockSortByPrice       StockSortBy = "PRICE"
	StockSortByStatus      StockSortBy = "STATUS"
)

func (s StockSortBy) IsValid() bool {
	switch s {
	case StockSortByStockedInAt, StockSortByMake, StockSortByModel,
		StockSortByModelYear, StockSortByPrice, StockSortByStatus:
		return true
	default:
		return false
	}
}

type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

func (s SortOrder) IsValid() bool { return s == SortOrderAsc || s == SortOrderDesc }

type InventoryStockFilters struct {
	Search     string
	Make       string
	Model      string
	Status     constant.StockStatus
	MinAgeDays *int
	MaxAgeDays *int
	AgingOnly  bool
	SortBy     StockSortBy
	SortOrder  SortOrder
	Page       int
	PageSize   int
}

type InventoryActionAggregate struct {
	ID         string                       `json:"id"`
	ActionType constant.InventoryActionType `json:"actionType"`
	Note       string                       `json:"note"`
	CreatedAt  time.Time                    `json:"createdAt"`
}

func (a *InventoryActionAggregate) FromModel(m *model.InventoryAction) {
	if a == nil || m == nil {
		return
	}
	a.ID = m.ID
	a.ActionType = constant.InventoryActionType(m.ActionType)
	a.Note = m.Note
	a.CreatedAt = m.CreatedAt
}

type InventoryStockAggregate struct {
	ID               string                    `json:"id"`
	DealershipID     string                    `json:"dealershipId"`
	VehicleID        string                    `json:"vehicleId"`
	VIN              string                    `json:"vin"`
	Make             string                    `json:"make"`
	Model            string                    `json:"model"`
	ModelYear        int                       `json:"modelYear"`
	Price            float64                   `json:"price"`
	Status           constant.StockStatus      `json:"status"`
	StockedInAt      time.Time                 `json:"stockedInAt"`
	StockedOutAt     *time.Time                `json:"stockedOutAt,omitempty"`
	InventoryAgeDays int                       `json:"inventoryAgeDays"`
	IsAging          bool                      `json:"isAging"`
	LatestAction     *InventoryActionAggregate `json:"latestAction,omitempty"`
}

func (a *InventoryStockAggregate) FromModel(stock *model.InventoryStock, action *model.InventoryAction, now time.Time) {
	if a == nil || stock == nil {
		return
	}
	age := 0
	status := constant.StockStatus(stock.Status)
	if status == constant.StockStatusInStock {
		age = calendarDaysBetween(stock.StockedInAt, now)
	}
	*a = InventoryStockAggregate{
		ID: stock.ID, DealershipID: stock.DealershipID, VehicleID: stock.VehicleID,
		VIN: stock.Vehicle.VIN, Make: stock.Vehicle.Make, Model: stock.Vehicle.Model,
		ModelYear: stock.Vehicle.ModelYear, Price: stock.Price, Status: status,
		StockedInAt: stock.StockedInAt, StockedOutAt: stock.StockedOutAt,
		InventoryAgeDays: age, IsAging: status == constant.StockStatusInStock && age > 90,
	}
	if action != nil {
		a.LatestAction = &InventoryActionAggregate{}
		a.LatestAction.FromModel(action)
	}
}

func calendarDaysBetween(from, to time.Time) int {
	fromDate := time.Date(from.UTC().Year(), from.UTC().Month(), from.UTC().Day(), 0, 0, 0, 0, time.UTC)
	toDate := time.Date(to.UTC().Year(), to.UTC().Month(), to.UTC().Day(), 0, 0, 0, 0, time.UTC)
	return int(toDate.Sub(fromDate).Hours() / 24)
}

type InventoryStockListAggregate struct {
	Items    []InventoryStockAggregate `json:"items"`
	Total    int64                     `json:"total"`
	Page     int                       `json:"page"`
	PageSize int                       `json:"pageSize"`
}

type CreateInventoryActionReq struct {
	ActionType constant.InventoryActionType `json:"actionType" validate:"required" enums:"PRICE_REDUCTION_PLANNED,TRANSFER_PROPOSED,MARKETING_CAMPAIGN,AWAITING_REVIEW,OTHER" example:"AWAITING_REVIEW"`
	Note       string                       `json:"note" validate:"required,min=1,max=500" minLength:"1" maxLength:"500" example:"Review this stock item with the sales manager"`
}

type CreateStockMovementReq struct {
	MovementType constant.StockMovementType `json:"movementType" validate:"required" enums:"STOCK_IN,STOCK_OUT" example:"STOCK_OUT"`
	Note         string                     `json:"note" validate:"required,min=1,max=500" minLength:"1" maxLength:"500" example:"Vehicle delivered to customer"`
}
