package inventory

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/service"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) CreateManualInventoryOperation(ctx context.Context, params models.CreateInventoryOperationParams) (*models.InventoryOperation, error) {
	if strings.TrimSpace(params.IngredientID) == "" || params.Amount <= 0 {
		return nil, utils.ValidationError("ingredient_id and positive amount are required", nil)
	}
	var op *models.InventoryOperation
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		activeShift, err := s.shiftRepo.GetActiveShift(ctx, true)
		if err != nil {
			return err
		}
		if activeShift == nil {
			return utils.BusinessError(utils.ErrorCodeShiftNotActive, "Нельзя создать складскую операцию без активной смены", nil)
		}
		ingredients, err := s.ingredientRepo.GetIngredientsByIDs(ctx, []string{params.IngredientID}, true)
		if err != nil {
			return err
		}
		if len(ingredients) == 0 {
			return utils.NotFoundError("Ингредиент не найден")
		}
		ingredient := ingredients[0]
		changeAmount := params.Amount
		if params.OperationType == service.OperationTypeManualWriteOff {
			if ingredient.CurrentStock+1e-9 < params.Amount {
				return utils.BusinessError(utils.ErrorCodeInsufficientStock, "Недостаточно остатка для списания", map[string]any{"ingredient_id": ingredient.ID, "ingredient_name": ingredient.Name, "requested": params.Amount, "available": ingredient.CurrentStock, "unit": ingredient.Unit})
			}
			changeAmount = -params.Amount
		}
		if err := s.ingredientRepo.AdjustIngredientStock(ctx, ingredient.ID, changeAmount); err != nil {
			return err
		}
		op, err = s.inventoryRepo.CreateInventoryOperation(ctx, models.InventoryOperationCreateRecord{ShiftID: activeShift.ID, IngredientID: ingredient.ID, ChangeAmount: changeAmount, OperationType: params.OperationType})
		return err
	})
	if err != nil {
		return nil, err
	}
	return op, nil
}
