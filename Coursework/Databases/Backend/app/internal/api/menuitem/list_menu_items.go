package menuitem

import (
	"strconv"

	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// listMenuItems godoc
// @Summary List menu items
// @Description Returns menu items with recipes and optional inactive records.
// @Tags menu-item
// @Produce json
// @Param include_inactive query bool false "Include inactive menu items" default(true)
// @Success 200 {object} apimodels.ListMenuItemsResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/menu-item [get]
func (a *API) listMenuItems(ctx *gin.Context) error {
	includeInactive := true
	if raw := ctx.Query("include_inactive"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return utils.ValidationError("invalid query parameter", map[string]any{"field": "include_inactive"})
		}
		includeInactive = value
	}
	items, err := a.service.ListMenuItems(ctx.Request.Context(), includeInactive)
	if err != nil {
		return err
	}
	responseItems := make([]apimodels.MenuItem, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, toMenuItemResponse(item))
	}
	ctx.JSON(200, apimodels.ListMenuItemsResponse{MenuItems: responseItems})
	return nil
}
