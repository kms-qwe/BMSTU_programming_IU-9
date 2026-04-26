package order

import (
	"strings"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createOrder godoc
// @Summary Create order
// @Description Creates an order, writes off ingredients and creates inventory operations in one transaction.
// @Tags order
// @Accept json
// @Produce json
// @Param request body apimodels.CreateOrderRequest true "Create order payload"
// @Success 200 {object} apimodels.CreateOrderResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/order [post]
func (a *API) createOrder(ctx *gin.Context) error {
	var request apimodels.CreateOrderRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	if len(request.Items) == 0 {
		return utils.ValidationError("items are required", nil)
	}
	items := make([]models.CreateOrderItemInput, 0, len(request.Items))
	for _, item := range request.Items {
		if strings.TrimSpace(item.MenuItemID) == "" || item.Quantity <= 0 {
			return utils.ValidationError("each item must contain menu_item_id and positive quantity", nil)
		}
		items = append(items, models.CreateOrderItemInput{MenuItemID: item.MenuItemID, Quantity: item.Quantity})
	}
	order, err := a.service.CreateOrder(ctx.Request.Context(), models.CreateOrderParams{Items: items})
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.CreateOrderResponse{Order: toOrderResponsePtr(order)})
	return nil
}
