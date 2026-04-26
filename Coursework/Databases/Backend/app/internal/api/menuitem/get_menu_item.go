package menuitem

import (
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// getMenuItem godoc
// @Summary Get menu item
// @Description Returns one menu item with its recipe.
// @Tags menu-item
// @Produce json
// @Param menu_item_id path string true "Menu item ID"
// @Success 200 {object} apimodels.GetMenuItemResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/menu-item/{menu_item_id} [get]
func (a *API) getMenuItem(ctx *gin.Context) error {
	item, err := a.service.GetMenuItem(ctx.Request.Context(), ctx.Param("menu_item_id"))
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.GetMenuItemResponse{MenuItem: toMenuItemResponsePtr(item)})
	return nil
}
