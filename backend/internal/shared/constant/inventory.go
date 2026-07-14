package constant

type StockStatus string

const (
	StockStatusInStock    StockStatus = "IN_STOCK"
	StockStatusOutOfStock StockStatus = "OUT_OF_STOCK"
)

func (s StockStatus) IsValid() bool {
	return s == StockStatusInStock || s == StockStatusOutOfStock
}

type StockMovementType string

const (
	StockMovementTypeIn  StockMovementType = "STOCK_IN"
	StockMovementTypeOut StockMovementType = "STOCK_OUT"
)

func (t StockMovementType) IsValid() bool {
	return t == StockMovementTypeIn || t == StockMovementTypeOut
}

type InventoryActionType string

const (
	InventoryActionTypePriceReductionPlanned InventoryActionType = "PRICE_REDUCTION_PLANNED"
	InventoryActionTypeTransferProposed      InventoryActionType = "TRANSFER_PROPOSED"
	InventoryActionTypeMarketingCampaign     InventoryActionType = "MARKETING_CAMPAIGN"
	InventoryActionTypeAwaitingReview        InventoryActionType = "AWAITING_REVIEW"
	InventoryActionTypeOther                 InventoryActionType = "OTHER"
)

func (t InventoryActionType) IsValid() bool {
	switch t {
	case InventoryActionTypePriceReductionPlanned,
		InventoryActionTypeTransferProposed,
		InventoryActionTypeMarketingCampaign,
		InventoryActionTypeAwaitingReview,
		InventoryActionTypeOther:
		return true
	default:
		return false
	}
}
