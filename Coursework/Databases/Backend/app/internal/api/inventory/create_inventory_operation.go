package inventory

import (
	"strings"

	"coffee-shop-backend/app/internal/api/common"
	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/utils"
	"coffee-shop-backend/app/pkg/apimodels"

	"github.com/gin-gonic/gin"
)

// createInventoryOperation godoc
// @Summary Create inventory operation
// @Description Creates a manual restock or manual write-off operation.
// @Tags inventory
// @Accept json
// @Produce json
// @Param request body apimodels.CreateOperationRequest true "Create inventory operation payload"
// @Success 200 {object} apimodels.CreateOperationResponse
// @Failure 400 {object} apimodels.ErrorResponse
// @Failure 404 {object} apimodels.ErrorResponse
// @Failure 500 {object} apimodels.ErrorResponse
// @Router /api/inventory/operation [post]
func (a *API) createInventoryOperation(ctx *gin.Context) error {
	var request apimodels.CreateOperationRequest
	if err := common.DecodeBody(ctx, &request); err != nil {
		return err
	}
	if strings.TrimSpace(request.IngredientID) == "" || strings.TrimSpace(request.OperationType) == "" || request.Amount <= 0 {
		return utils.ValidationError("ingredient_id, operation_type and positive amount are required", nil)
	}
	item, err := a.service.CreateManualInventoryOperation(ctx.Request.Context(), models.CreateInventoryOperationParams{
		IngredientID: request.IngredientID, OperationType: request.OperationType, Amount: request.Amount,
	})
	if err != nil {
		return err
	}
	ctx.JSON(200, apimodels.CreateOperationResponse{Operation: toInventoryOperationResponsePtr(item)})
	return nil
}
