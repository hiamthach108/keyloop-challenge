package handler

import (
	"strconv"

	"github.com/hiamthach108/dreon-sdk/errorx"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/aggregate"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/service"
	"github.com/hiamthach108/keyloop-challenge/backend/internal/shared/constant"
	"github.com/labstack/echo/v4"
)

type InventoryHandler struct{ service service.IInventorySvc }

func NewInventoryHandler(service service.IInventorySvc) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) RegisterRoutes(group *echo.Group) {
	group.GET("", h.HandleListDealerships)
	group.GET("/:dealershipID/stocks", h.HandleListStocks)
	group.GET("/:dealershipID/stocks/aging", h.HandleListAgingStocks)
	group.POST("/:dealershipID/stocks/:stockID/actions", h.HandleCreateAction)
	group.POST("/:dealershipID/stocks/:stockID/movements", h.HandleRecordMovement)
	group.GET("/:dealershipID/stocks/:stockID/history", h.HandleGetHistory)
}

func (h *InventoryHandler) HandleListDealerships(c echo.Context) error {
	items, err := h.service.ListDealerships(c.Request().Context())
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, items)
}

func (h *InventoryHandler) HandleListStocks(c echo.Context) error {
	filters, err := parseStockFilters(c)
	if err != nil {
		return HandleError(c, err)
	}
	items, err := h.service.ListStocks(c.Request().Context(), c.Param("dealershipID"), filters)
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, items)
}

func (h *InventoryHandler) HandleListAgingStocks(c echo.Context) error {
	filters, err := parseStockFilters(c)
	if err != nil {
		return HandleError(c, err)
	}
	filters.AgingOnly = true
	filters.Status = constant.StockStatusInStock
	items, err := h.service.ListStocks(c.Request().Context(), c.Param("dealershipID"), filters)
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, items)
}

func (h *InventoryHandler) HandleCreateAction(c echo.Context) error {
	req, err := HandleValidateBind[aggregate.CreateInventoryActionReq](c)
	if err != nil {
		return HandleError(c, err)
	}
	action, err := h.service.CreateAction(c.Request().Context(), c.Param("dealershipID"), c.Param("stockID"), &req)
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, action)
}

func (h *InventoryHandler) HandleRecordMovement(c echo.Context) error {
	req, err := HandleValidateBind[aggregate.CreateStockMovementReq](c)
	if err != nil {
		return HandleError(c, err)
	}
	stock, err := h.service.RecordMovement(c.Request().Context(), c.Param("dealershipID"), c.Param("stockID"), &req)
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, stock)
}

func (h *InventoryHandler) HandleGetHistory(c echo.Context) error {
	history, err := h.service.GetHistory(c.Request().Context(), c.Param("dealershipID"), c.Param("stockID"))
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, history)
}

func parseStockFilters(c echo.Context) (aggregate.InventoryStockFilters, error) {
	filters := aggregate.InventoryStockFilters{
		Search: c.QueryParam("search"), Make: c.QueryParam("make"), Model: c.QueryParam("model"),
		Status:    constant.StockStatus(c.QueryParam("status")),
		SortBy:    aggregate.StockSortBy(c.QueryParam("sortBy")),
		SortOrder: aggregate.SortOrder(c.QueryParam("sortOrder")),
	}
	var err error
	if filters.MinAgeDays, err = optionalIntQuery(c, "minAgeDays"); err != nil {
		return filters, err
	}
	if filters.MaxAgeDays, err = optionalIntQuery(c, "maxAgeDays"); err != nil {
		return filters, err
	}
	if value := c.QueryParam("agingOnly"); value != "" {
		filters.AgingOnly, err = strconv.ParseBool(value)
		if err != nil {
			return filters, errorx.New(errorx.ErrBadRequest, "agingOnly must be a boolean")
		}
	}
	if filters.Page, err = intQueryWithDefault(c, "page", 1); err != nil {
		return filters, err
	}
	if filters.PageSize, err = intQueryWithDefault(c, "pageSize", 20); err != nil {
		return filters, err
	}
	return filters, nil
}

func optionalIntQuery(c echo.Context, name string) (*int, error) {
	value := c.QueryParam(name)
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil, errorx.New(errorx.ErrBadRequest, name+" must be an integer")
	}
	return &parsed, nil
}

func intQueryWithDefault(c echo.Context, name string, defaultValue int) (int, error) {
	value := c.QueryParam(name)
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errorx.New(errorx.ErrBadRequest, name+" must be an integer")
	}
	return parsed, nil
}
