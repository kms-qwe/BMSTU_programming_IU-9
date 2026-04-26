package ingredient

import (
	"context"
	"strings"

	"coffee-shop-backend/app/internal/models"
	"coffee-shop-backend/app/internal/service"
	"coffee-shop-backend/app/internal/utils"
)

func (s *Service) CreateIngredient(ctx context.Context, params models.CreateIngredientParams) (*models.Ingredient, error) {
	params.Name = strings.TrimSpace(params.Name)
	params.Unit = strings.TrimSpace(params.Unit)
	if params.Name == "" || params.Unit == "" {
		return nil, utils.ValidationError("name and unit are required", nil)
	}
	if params.InitialStock < 0 {
		return nil, utils.ValidationError("initial_stock must be non-negative", nil)
	}
	var ingredient *models.Ingredient
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var err error
		ingredient, err = s.ingredientRepo.CreateIngredient(ctx, params)
		if err != nil || params.InitialStock <= 0 {
			return err
		}
		activeShift, err := s.shiftRepo.GetActiveShift(ctx, true)
		if err != nil {
			return err
		}
		if activeShift == nil {
			return utils.BusinessError(utils.ErrorCodeShiftNotActive, "Нельзя создать начальный остаток без активной смены", nil)
		}
		_, err = s.inventoryRepo.CreateInventoryOperation(ctx, models.InventoryOperationCreateRecord{
			ShiftID: activeShift.ID, IngredientID: ingredient.ID, ChangeAmount: params.InitialStock, OperationType: service.OperationTypeInitialStock,
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	return ingredient, nil
}
