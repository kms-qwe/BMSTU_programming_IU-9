package menuitem

import (
	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// updateMenuItemPrice godoc
// @Summary Update menu item price
// @Description Updates only the price of an existing menu item.
// @Tags menu-item
// @Accept json
// @Produce json
// @Param menu_item_id path string true "Menu item ID"
// @Param request body apimodels.UpdateMenuItemPriceRequest true "Update price payload"
// @Success 200 {object} apimodels.UpdateMenuItemPriceResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/menu-item/{menu_item_id}/price [patch]
func (a *API) updateMenuItemPrice(ctx *gin.Context) error {
	var request apimodels.UpdateMenuItemPriceRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	if request.Price < 0 {
		return utils.ValidationError("price must be non-negative", nil)
	}
	item, err := a.service.UpdateMenuItemPrice(ctx.Request.Context(), models.UpdateMenuItemPriceParams{
		MenuItemID: ctx.Param("menu_item_id"),
		Price:      request.Price,
	})
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.UpdateMenuItemPriceResponse{MenuItem: toMenuItemResponsePtr(item)})
	return nil
}
