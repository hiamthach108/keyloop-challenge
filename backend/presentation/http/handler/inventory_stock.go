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

// HandleListDealerships godoc
// @Summary List dealerships
// @Description Returns seeded dealerships available to the inventory dashboard.
// @Tags dealerships
// @Produce json
// @Success 200 {object} DealershipListResp
// @Failure 500 {object} BaseResp
// @Router /dealerships [get]
func (h *InventoryHandler) HandleListDealerships(c echo.Context) error {
	items, err := h.service.ListDealerships(c.Request().Context())
	if err != nil {
		return HandleError(c, err)
	}
	return HandleSuccess(c, items)
}

// HandleListStocks godoc
// @Summary List dealership stock
// @Description Returns paginated inventory stock with search, filters, and enum-constrained sorting.
// @Tags inventory
// @Produce json
// @Param dealershipID path string true "Dealership ID"
// @Param search query string false "Case-insensitive make/model search"
// @Param make query string false "Case-insensitive make filter"
// @Param model query string false "Case-insensitive model filter"
// @Param status query string false "Stock status" Enums(IN_STOCK, OUT_OF_STOCK)
// @Param minAgeDays query int false "Minimum stock age in calendar days" minimum(0)
// @Param maxAgeDays query int false "Maximum stock age in calendar days" minimum(0)
// @Param agingOnly query bool false "Only return in-stock inventory older than 90 days"
// @Param sortBy query string false "Sort column" Enums(STOCKED_IN_AT, MAKE, MODEL, MODEL_YEAR, PRICE, STATUS)
// @Param sortOrder query string false "Sort order" Enums(ASC, DESC)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param pageSize query int false "Page size" default(20) minimum(1) maximum(100)
// @Success 200 {object} StockListResp
// @Failure 400 {object} ValidationErrResp
// @Failure 404 {object} BaseResp
// @Failure 500 {object} BaseResp
// @Router /dealerships/{dealershipID}/stocks [get]
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

// HandleListAgingStocks godoc
// @Summary List aging stock
// @Description Returns only in-stock inventory older than 90 calendar days.
// @Tags inventory
// @Produce json
// @Param dealershipID path string true "Dealership ID"
// @Param search query string false "Case-insensitive make/model search"
// @Param make query string false "Case-insensitive make filter"
// @Param model query string false "Case-insensitive model filter"
// @Param sortBy query string false "Sort column" Enums(STOCKED_IN_AT, MAKE, MODEL, MODEL_YEAR, PRICE, STATUS)
// @Param sortOrder query string false "Sort order" Enums(ASC, DESC)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param pageSize query int false "Page size" default(20) minimum(1) maximum(100)
// @Success 200 {object} StockListResp
// @Failure 400 {object} ValidationErrResp
// @Failure 404 {object} BaseResp
// @Failure 500 {object} BaseResp
// @Router /dealerships/{dealershipID}/stocks/aging [get]
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

// HandleCreateAction godoc
// @Summary Create stock action
// @Description Appends a manager action for an in-stock vehicle older than 90 days.
// @Tags actions
// @Accept json
// @Produce json
// @Param dealershipID path string true "Dealership ID"
// @Param stockID path string true "Stock ID"
// @Param request body aggregate.CreateInventoryActionReq true "Action payload"
// @Success 200 {object} InventoryActionResp
// @Failure 400 {object} ValidationErrResp
// @Failure 404 {object} BaseResp
// @Failure 422 {object} BaseResp
// @Failure 500 {object} BaseResp
// @Router /dealerships/{dealershipID}/stocks/{stockID}/actions [post]
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

// HandleRecordMovement godoc
// @Summary Record stock movement
// @Description Atomically appends a STOCK_IN or STOCK_OUT movement and transitions current stock state.
// @Tags movements
// @Accept json
// @Produce json
// @Param dealershipID path string true "Dealership ID"
// @Param stockID path string true "Stock ID"
// @Param request body aggregate.CreateStockMovementReq true "Movement payload"
// @Success 200 {object} StockResp
// @Failure 400 {object} ValidationErrResp
// @Failure 404 {object} BaseResp
// @Failure 409 {object} BaseResp
// @Failure 500 {object} BaseResp
// @Router /dealerships/{dealershipID}/stocks/{stockID}/movements [post]
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

// HandleGetHistory godoc
// @Summary Get stock history
// @Description Returns merged stock movement and manager action history newest first.
// @Tags inventory
// @Produce json
// @Param dealershipID path string true "Dealership ID"
// @Param stockID path string true "Stock ID"
// @Success 200 {object} StockHistoryResp
// @Failure 400 {object} ValidationErrResp
// @Failure 404 {object} BaseResp
// @Failure 500 {object} BaseResp
// @Router /dealerships/{dealershipID}/stocks/{stockID}/history [get]
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
