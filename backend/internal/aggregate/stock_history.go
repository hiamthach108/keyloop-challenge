package aggregate

import (
	"sort"
	"time"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/shared/constant"
)

type StockHistoryEventType string

const (
	StockHistoryEventTypeAction   StockHistoryEventType = "ACTION"
	StockHistoryEventTypeStockIn  StockHistoryEventType = "STOCK_IN"
	StockHistoryEventTypeStockOut StockHistoryEventType = "STOCK_OUT"
)

type StockHistoryEventAggregate struct {
	ID         string                        `json:"id"`
	EventType  StockHistoryEventType         `json:"eventType"`
	ActionType *constant.InventoryActionType `json:"actionType,omitempty"`
	Note       string                        `json:"note"`
	OccurredAt time.Time                     `json:"occurredAt"`
}

func StockHistoryFromModels(actions []model.InventoryAction, movements []model.StockMovement) []StockHistoryEventAggregate {
	result := make([]StockHistoryEventAggregate, 0, len(actions)+len(movements))
	for _, action := range actions {
		actionType := constant.InventoryActionType(action.ActionType)
		result = append(result, StockHistoryEventAggregate{
			ID: action.ID, EventType: StockHistoryEventTypeAction, ActionType: &actionType,
			Note: action.Note, OccurredAt: action.CreatedAt,
		})
	}
	for _, movement := range movements {
		eventType := StockHistoryEventTypeStockIn
		if constant.StockMovementType(movement.MovementType) == constant.StockMovementTypeOut {
			eventType = StockHistoryEventTypeStockOut
		}
		result = append(result, StockHistoryEventAggregate{
			ID: movement.ID, EventType: eventType, Note: movement.Note, OccurredAt: movement.OccurredAt,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].OccurredAt.After(result[j].OccurredAt) })
	return result
}
