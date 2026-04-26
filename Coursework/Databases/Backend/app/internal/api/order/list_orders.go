package order

import (
	"time"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// listOrders godoc
// @Summary List orders
// @Description Returns orders with employee, shift and item details.
// @Tags order
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param employee_id query string false "Employee ID"
// @Param from query string false "RFC3339 from timestamp"
// @Param to query string false "RFC3339 to timestamp"
// @Success 200 {object} apimodels.ListOrdersResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/order [get]
func (a *API) listOrders(ctx *gin.Context) error {
	page, pageSize, err := common.ParsePagination(ctx, 1, 20)
	if err != nil {
		return err
	}
	var filter models.OrderFilter
	filter.Page, filter.PageSize = page, pageSize
	filter.EmployeeID = ctx.Query("employee_id")
	if raw := ctx.Query("from"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return utils.ValidationError("invalid from date", nil)
		}
		filter.From = &t
	}
	if raw := ctx.Query("to"); raw != "" {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			return utils.ValidationError("invalid to date", nil)
		}
		filter.To = &t
	}
	result, err := a.service.ListOrders(ctx.Request.Context(), filter)
	if err != nil {
		return err
	}
	items := make([]apimodels.Order, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toOrderResponse(item))
	}
	ctx.JSON(200, apimodels.ListOrdersResponse{Items: items, Pagination: apimodels.Pagination{
		Page: result.Pagination.Page, PageSize: result.Pagination.PageSize, TotalItems: result.Pagination.TotalItems, TotalPages: result.Pagination.TotalPages,
	}})
	return nil
}
