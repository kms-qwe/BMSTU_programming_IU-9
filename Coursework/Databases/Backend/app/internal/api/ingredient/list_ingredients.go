package ingredient

import (
	"strconv"

	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// listIngredients godoc
// @Summary List ingredients
// @Description Returns all ingredients with optional inactive records.
// @Tags ingredient
// @Produce json
// @Param include_inactive query bool false "Include inactive ingredients" default(true)
// @Success 200 {object} apimodels.ListIngredientsResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/ingredient [get]
func (a *API) listIngredients(ctx *gin.Context) error {
	includeInactive := true
	if raw := ctx.Query("include_inactive"); raw != "" {
		value, err := strconv.ParseBool(raw)
		if err != nil {
			return utils.ValidationError("invalid query parameter", map[string]any{"field": "include_inactive"})
		}
		includeInactive = value
	}

	items, err := a.service.ListIngredients(ctx.Request.Context(), includeInactive)
	if err != nil {
		return err
	}

	responseItems := make([]apimodels.Ingredient, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, toIngredientResponse(item))
	}

	ctx.JSON(200, apimodels.ListIngredientsResponse{Ingredients: responseItems})
	return nil
}
