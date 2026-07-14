package aggregate

import (
	"testing"
	"time"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
)

func TestStockHistoryFromModelsSortsAllEventsNewestFirst(t *testing.T) {
	t.Parallel()
	base := time.Date(2026, time.July, 12, 0, 0, 0, 0, time.UTC)
	actions := []model.InventoryAction{{ID: "action", ActionType: "OTHER", CreatedAt: base.Add(time.Hour)}}
	movements := []model.StockMovement{
		{ID: "in", MovementType: "STOCK_IN", OccurredAt: base},
		{ID: "out", MovementType: "STOCK_OUT", OccurredAt: base.Add(2 * time.Hour)},
	}

	got := StockHistoryFromModels(actions, movements)

	if len(got) != 3 {
		t.Fatalf("history length = %d, want 3", len(got))
	}
	wantIDs := []string{"out", "action", "in"}
	for i, wantID := range wantIDs {
		if got[i].ID != wantID {
			t.Fatalf("history[%d].ID = %q, want %q", i, got[i].ID, wantID)
		}
	}
	if got[0].EventType != StockHistoryEventTypeStockOut || got[1].EventType != StockHistoryEventTypeAction {
		t.Fatalf("unexpected event types: %q, %q", got[0].EventType, got[1].EventType)
	}
}
