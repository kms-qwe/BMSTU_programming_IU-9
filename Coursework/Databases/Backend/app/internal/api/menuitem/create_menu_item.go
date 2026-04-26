package menuitem

import (
	"strings"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createMenuItem godoc
// @Summary Create menu item
// @Description Creates a new menu item with fixed recipe.
// @Tags menu-item
// @Accept json
// @Produce json
// @Param request body apimodels.CreateMenuItemRequest true "Create menu item payload"
// @Success 200 {object} apimodels.CreateMenuItemResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 409 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/menu-item [post]
func (a *API) createMenuItem(ctx *gin.Context) error {
	var request apimodels.CreateMenuItemRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	if strings.TrimSpace(request.Name) == "" || request.Price < 0 || len(request.Recipe) == 0 {
		return utils.ValidationError("invalid menu item payload", nil)
	}
	recipe := make([]models.CreateMenuItemRecipeInput, 0, len(request.Recipe))
	for _, item := range request.Recipe {
		if strings.TrimSpace(item.IngredientID) == "" || item.Amount <= 0 {
			return utils.ValidationError("recipe items must contain ingredient_id and positive amount", nil)
		}
		recipe = append(recipe, models.CreateMenuItemRecipeInput{IngredientID: item.IngredientID, Amount: item.Amount})
	}
	menuItem, err := a.service.CreateMenuItem(ctx.Request.Context(), models.CreateMenuItemParams{Name: request.Name, Price: request.Price, Recipe: recipe})
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.CreateMenuItemResponse{MenuItem: toMenuItemResponsePtr(menuItem)})
	return nil
}
