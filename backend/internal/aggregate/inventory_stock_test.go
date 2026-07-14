package aggregate

import (
	"testing"
	"time"

	"github.com/hiamthach108/keyloop-challenge/backend/internal/model"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/shared/constant"
)

func TestInventoryStockAggregateAgingBoundary(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, time.July, 12, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name      string
		status    constant.StockStatus
		stockedAt time.Time
		wantAge   int
		wantAging bool
	}{
		{name: "90 days is not aging", status: constant.StockStatusInStock, stockedAt: now.AddDate(0, 0, -90), wantAge: 90},
		{name: "91 days is aging", status: constant.StockStatusInStock, stockedAt: now.AddDate(0, 0, -91), wantAge: 91, wantAging: true},
		{name: "out of stock is never aging", status: constant.StockStatusOutOfStock, stockedAt: now.AddDate(0, 0, -120)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stock := model.InventoryStock{Status: string(tt.status), StockedInAt: tt.stockedAt}
			var got InventoryStockAggregate
			got.FromModel(&stock, nil, now)
			if got.InventoryAgeDays != tt.wantAge {
				t.Fatalf("InventoryAgeDays = %d, want %d", got.InventoryAgeDays, tt.wantAge)
			}
			if got.IsAging != tt.wantAging {
				t.Fatalf("IsAging = %t, want %t", got.IsAging, tt.wantAging)
			}
		})
	}
}
