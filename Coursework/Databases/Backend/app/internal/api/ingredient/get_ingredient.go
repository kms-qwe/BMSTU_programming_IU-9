package ingredient

import (
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// getIngredient godoc
// @Summary Get ingredient
// @Description Returns one ingredient with recipe usage details.
// @Tags ingredient
// @Produce json
// @Param ingredient_id path string true "Ingredient ID"
// @Success 200 {object} apimodels.GetIngredientResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/ingredient/{ingredient_id} [get]
func (a *API) getIngredient(ctx *gin.Context) error {
	item, err := a.service.GetIngredient(ctx.Request.Context(), ctx.Param("ingredient_id"))
	if err != nil {
		return err
	}

	ctx.JSON(200, apimodels.GetIngredientResponse{Ingredient: toIngredientDetailsResponse(item)})
	return nil
}
