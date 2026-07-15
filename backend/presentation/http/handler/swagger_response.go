package handler

import "github.com/hiamthach108/keyloop-challenge/backend/internal/aggregate"

type PingResp struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"pong"`
}

type DealershipListResp struct {
	Code    int                             `json:"code" example:"200"`
	Message string                          `json:"message" example:"success"`
	Data    []aggregate.DealershipAggregate `json:"data"`
}

type StockListResp struct {
	Code    int                                   `json:"code" example:"200"`
	Message string                                `json:"message" example:"success"`
	Data    aggregate.InventoryStockListAggregate `json:"data"`
}

type StockResp struct {
	Code    int                               `json:"code" example:"200"`
	Message string                            `json:"message" example:"success"`
	Data    aggregate.InventoryStockAggregate `json:"data"`
}

type InventoryActionResp struct {
	Code    int                                `json:"code" example:"200"`
	Message string                             `json:"message" example:"success"`
	Data    aggregate.InventoryActionAggregate `json:"data"`
}

type StockHistoryResp struct {
	Code    int                                    `json:"code" example:"200"`
	Message string                                 `json:"message" example:"success"`
	Data    []aggregate.StockHistoryEventAggregate `json:"data"`
}
