package ingredient

import (
	"strings"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createIngredient godoc
// @Summary Create ingredient
// @Description Creates a new ingredient and optionally initializes stock.
// @Tags ingredient
// @Accept json
// @Produce json
// @Param request body apimodels.CreateIngredientRequest true "Create ingredient payload"
// @Success 200 {object} apimodels.CreateIngredientResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 409 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/ingredient [post]
func (a *API) createIngredient(ctx *gin.Context) error {
	var request apimodels.CreateIngredientRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	if strings.TrimSpace(request.Name) == "" || strings.TrimSpace(request.Unit) == "" {
		return utils.ValidationError("name and unit are required", nil)
	}
	if request.InitialStock != nil && *request.InitialStock < 0 {
		return utils.ValidationError("initial_stock must be non-negative", nil)
	}

	initialStock := 0.0
	if request.InitialStock != nil {
		initialStock = *request.InitialStock
	}

	item, err := a.service.CreateIngredient(ctx.Request.Context(), models.CreateIngredientParams{
		Name:         request.Name,
		Unit:         request.Unit,
		InitialStock: initialStock,
	})
	if err != nil {
		return err
	}

	ctx.JSON(200, apimodels.CreateIngredientResponse{Ingredient: toIngredientResponse(*item)})
	return nil
}
