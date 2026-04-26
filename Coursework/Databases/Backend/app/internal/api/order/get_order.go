package order

import (
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// getOrder godoc
// @Summary Get order
// @Description Returns one order with items and linked inventory operations.
// @Tags order
// @Produce json
// @Param order_id path string true "Order ID"
// @Success 200 {object} apimodels.GetOrderResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/order/{order_id} [get]
func (a *API) getOrder(ctx *gin.Context) error {
	item, err := a.service.GetOrder(ctx.Request.Context(), ctx.Param("order_id"))
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.GetOrderResponse{Order: toOrderResponsePtr(item)})
	return nil
}
