package inventory

import (
	"time"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// listInventoryOperations godoc
// @Summary List inventory operations
// @Description Returns inventory operations with filters by employee, ingredient, type and date range.
// @Tags inventory
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param employee_id query string false "Employee ID"
// @Param ingredient_id query string false "Ingredient ID"
// @Param operation_type query string false "Operation type"
// @Param from query string false "RFC3339 from timestamp"
// @Param to query string false "RFC3339 to timestamp"
// @Success 200 {object} apimodels.ListOperationsResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/inventory/operation [get]
func (a *API) listInventoryOperations(ctx *gin.Context) error {
	page, pageSize, err := common.ParsePagination(ctx, 1, 20)
	if err != nil {
		return err
	}
	filter := models.InventoryOperationFilter{Page: page, PageSize: pageSize, EmployeeID: ctx.Query("employee_id"), IngredientID: ctx.Query("ingredient_id"), OperationType: ctx.Query("operation_type")}
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
	result, err := a.service.ListInventoryOperations(ctx.Request.Context(), filter)
	if err != nil {
		return err
	}
	items := make([]apimodels.InventoryOperation, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toInventoryOperationResponse(item))
	}
	ctx.JSON(200, apimodels.ListOperationsResponse{Items: items, Pagination: apimodels.Pagination{
		Page: result.Pagination.Page, PageSize: result.Pagination.PageSize, TotalItems: result.Pagination.TotalItems, TotalPages: result.Pagination.TotalPages,
	}})
	return nil
}
