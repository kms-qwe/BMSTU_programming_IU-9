package menuitem

import (
	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// setMenuItemActivity godoc
// @Summary Set menu item activity
// @Description Activates or deactivates a menu item.
// @Tags menu-item
// @Accept json
// @Produce json
// @Param menu_item_id path string true "Menu item ID"
// @Param request body apimodels.SetMenuItemActivityRequest true "Activity payload"
// @Success 200 {object} apimodels.SetMenuItemActivityResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/menu-item/{menu_item_id}/activity [patch]
func (a *API) setMenuItemActivity(ctx *gin.Context) error {
	var request apimodels.SetMenuItemActivityRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	item, err := a.service.SetMenuItemActivity(ctx.Request.Context(), models.SetMenuItemActivityParams{
		MenuItemID: ctx.Param("menu_item_id"),
		IsActive:   request.IsActive,
	})
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.SetMenuItemActivityResponse{MenuItem: toMenuItemResponsePtr(item)})
	return nil
}
