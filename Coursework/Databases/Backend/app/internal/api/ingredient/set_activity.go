package ingredient

import (
	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// setIngredientActivity godoc
// @Summary Set ingredient activity
// @Description Activates or deactivates an ingredient.
// @Tags ingredient
// @Accept json
// @Produce json
// @Param ingredient_id path string true "Ingredient ID"
// @Param request body apimodels.SetIngredientActivityRequest true "Activity payload"
// @Success 200 {object} apimodels.SetIngredientActivityResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/ingredient/{ingredient_id}/activity [patch]
func (a *API) setIngredientActivity(ctx *gin.Context) error {
	var request apimodels.SetIngredientActivityRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}

	item, err := a.service.SetIngredientActivity(ctx.Request.Context(), models.SetIngredientActivityParams{
		IngredientID: ctx.Param("ingredient_id"),
		IsActive:     request.IsActive,
	})
	if err != nil {
		return err
	}

	ctx.JSON(200, apimodels.SetIngredientActivityResponse{Ingredient: toIngredientResponse(*item)})
	return nil
}
